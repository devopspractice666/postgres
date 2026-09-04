package main

import (
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func InstrumentHandler(handler http.HandlerFunc, method, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ActiveConnections.Inc()
		defer ActiveConnections.Dec()

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		handler(ww, r)

		duration := time.Since(start).Seconds()
		RequestsTotal.WithLabelValues(method, path, http.StatusText(ww.statusCode)).Inc()
		RequestDuration.WithLabelValues(method, path).Observe(duration)
	})
}
