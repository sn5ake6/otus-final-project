package internalhttp

import (
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.statusCode = http.StatusOK

	return w.ResponseWriter.Write(data)
}

func loggingMiddleware(h http.HandlerFunc, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{w, 0}
		start := time.Now()
		h(writer, r)
		logger.LogHTTPRequest(r, writer.statusCode, time.Since(start))
	}
}
