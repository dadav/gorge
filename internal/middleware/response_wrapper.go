package middleware

import "net/http"

type ResponseWrapper struct {
	http.ResponseWriter
	written bool
}

func NewResponseWrapper(w http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{ResponseWriter: w}
}

func (w *ResponseWrapper) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWrapper) WriteHeader(statusCode int) {
	w.written = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWrapper) WasWritten() bool {
	return w.written
}
