package server

import (
	"io"
	"log/slog"
	"maps"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"

	"github.com/ljgago/api-viewer/pkg/log"
)

func NewProxy(wg *sync.WaitGroup, opts Options) {
	defer wg.Done()

	// Start the proxy server to the last
	time.Sleep(100 * time.Millisecond)

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Recoverer)
	r.Use(httplog.RequestLogger(log.Get()))

	// Handlers
	r.HandleFunc("/*", handleScalarProxy)

	log.Infof("[Proxy server] listening on http://localhost:%s", opts.Proxy)
	if err := http.ListenAndServe("localhost:"+opts.Proxy, r); err != nil {
		log.Error("Error starting proxy server", slog.String("error", err.Error()))
	}
}

func handleScalarProxy(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Copy headers from the original request to maintain authentication
			// and other important headers through redirect chains.
			maps.Copy(req.Header, via[0].Header)

			return nil
		},
	}

	target := r.URL.Query().Get("scalar_url")

	// Create the outbound request
	outreq, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		log.Error("redirect to blocked", slog.String("host", r.URL.Host))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the headers but exclude Origin. It's not required and might confuse some target servers.
	for key, values := range r.Header {
		if !strings.EqualFold(strings.ToLower(key), "origin") {
			outreq.Header[key] = values
		}
	}

	// Redirect X-Scalar-Cookie header as Cookie header
	if xScalarCookie := r.Header.Get("X-Scalar-Cookie"); xScalarCookie != "" {
		// Set the cookie
		outreq.Header.Set("Cookie", xScalarCookie)
		// Remove the X-Scalar-Cookie header
		outreq.Header.Del("X-Scalar-Cookie")
	}

	// Make the request
	resp, err := client.Do(outreq)
	if err != nil {
		log.Error("make request", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Close response body when done to prevent resource leaks
	defer resp.Body.Close()

	// Copy headers from final response, but skip CORS headers
	for key, values := range resp.Header {
		// Check if header is a CORS headers
		isCorsHeader := func(header string) bool {
			return strings.HasPrefix(strings.ToLower(header), "access-control-")
		}

		if !isCorsHeader(key) {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	// Add CORS headers here, after the response headers are copied
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	// Add the final URL as a header
	w.Header().Set("X-Forwarded-Host", resp.Request.URL.String())

	// Copy the status code from the proxied response
	w.WriteHeader(resp.StatusCode)

	// Copy the body
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error("copy body", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}
