package log

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-chi/httplog/v2"
)

var logger *httplog.Logger

// Setup configurates the log instance.
//   - enableJSON bool - true | false
//   - level string - "debug" | "info" | "warn" | "error"
func Setup(enableJSON bool, level string) {
	logLevel := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	logger = httplog.NewLogger("api-viewer", httplog.Options{
		JSON:             enableJSON,
		LogLevel:         logLevel[level],
		Concise:          true,
		RequestHeaders:   true,
		ResponseHeaders:  true,
		MessageFieldName: "message",
		TimeFieldFormat:  time.RFC3339,
		TimeFieldName:    "time",
	})
}

func Get() *httplog.Logger {
	return logger
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Debugf(format string, a ...any) {
	logger.Debug(fmt.Sprintf(format, a...))
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Infof(format string, a ...any) {
	logger.Info(fmt.Sprintf(format, a...))
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Warnf(format string, a ...any) {
	logger.Warn(fmt.Sprintf(format, a...))
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Errorf(format string, a ...any) {
	logger.Error(fmt.Sprintf(format, a...))
}
