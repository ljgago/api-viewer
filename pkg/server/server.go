package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"

	"github.com/ljgago/api-viewer/assets"
	"github.com/ljgago/api-viewer/pkg/log"
	"github.com/ljgago/api-viewer/pkg/middleware/csrf"
)

type TemplateData struct {
	Url   string
	Proxy string
	Theme string
}

type Options struct {
	Port  string
	Proxy string
	File  string
	Theme string
}

func NewServer(wg *sync.WaitGroup, opts Options) {
	defer wg.Done()

	tmpl, err := template.ParseFS(assets.StaticFS, "index.html")
	if err != nil {
		log.Error("error parsing template", slog.String("error", err.Error()))
	}

	specname := filepath.Base(opts.File)
	specpath := filepath.Join("spec", specname)

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Recoverer)
	r.Use(httplog.RequestLogger(log.Get()))
	r.Use(csrf.Middleware)

	// Handlers
	r.HandleFunc("/", handleTemplate(tmpl, specpath, opts.Proxy, opts.Theme))
	r.HandleFunc("/"+specpath, handleSpec(opts.File))
	r.Handle("/*", http.FileServerFS(assets.StaticFS))

	log.Infof("[Main server] listening on http://localhost:%s", opts.Port)
	if err := http.ListenAndServe("localhost:"+opts.Port, r); err != nil {
		log.Error("Error starting server", slog.String("error", err.Error()))
	}
}

func handleTemplate(tmpl *template.Template, specpath, proxy, theme string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateData := TemplateData{
			Url:   fmt.Sprintf("/%s", specpath),
			Proxy: fmt.Sprintf("http://localhost:%s", proxy),
			Theme: theme,
		}

		if err := tmpl.Execute(w, templateData); err != nil {
			log.Error("template execute error", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func handleSpec(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}
