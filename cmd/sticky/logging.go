package main

import (
	"net/http"
	"os"

	"github.com/stickysh/sticky/pkg/id"
	kitlog "github.com/go-kit/kit/log"
)

func NewLogger() kitlog.Logger {
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "instance_id", id.GenULID())

	return logger
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{w, http.StatusOK}
}

func (lrw *responseWriterWrapper) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type loggerHandler struct {
	h http.Handler
	logger kitlog.Logger
}

func NewLoggerMiddelware(logger kitlog.Logger, handler http.Handler) http.Handler {
	return &loggerHandler{
		h: handler,
		logger: logger,
	}
}

func (l *loggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ww := NewResponseWriterWrapper(w)
	l.h.ServeHTTP(ww, r)
	l.logger.Log("method", r.Method, "uri", r.URL, "status", ww.statusCode)
}



