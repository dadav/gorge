package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/dadav/gorge/internal/log"
)

// capturedResponseWriter is a custom response writer that captures the response status
type capturedResponseWriter struct {
	http.ResponseWriter
	body   []byte
	status int
}

func (w *capturedResponseWriter) WriteHeader(code int) {
	w.status = code
}

func (w *capturedResponseWriter) Write(body []byte) (int, error) {
	w.body = body
	return len(body), nil
}

func (w *capturedResponseWriter) sendCapturedResponse() {
	w.ResponseWriter.WriteHeader(w.status)
	w.ResponseWriter.Write(w.body)
}

func ProxyFallback(upstreamHost string, forwardToProxy func(int) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// capture response
			capturedResponseWriter := &capturedResponseWriter{ResponseWriter: w}
			next.ServeHTTP(capturedResponseWriter, r)

			if forwardToProxy(capturedResponseWriter.status) {
				log.Log.Infof("Forwarding request to %s\n", upstreamHost)
				forwardRequest(w, r, upstreamHost)
				return
			}

			// If the response status is not 404, serve the original response
			capturedResponseWriter.sendCapturedResponse()
		})
	}
}

func forwardRequest(w http.ResponseWriter, r *http.Request, forwardHost string) {
	// Create a buffer to store the request body
	var requestBodyBytes []byte
	if r.Body != nil {
		requestBodyBytes, _ = io.ReadAll(r.Body)
	}

	// Clone the original request
	forwardUrl, err := url.JoinPath(forwardHost, r.URL.Path)
	if err != nil {
		http.Error(w, "Failed to create forwarded request", http.StatusInternalServerError)
		return
	}

	forwardedRequest, err := http.NewRequest(r.Method, forwardUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		http.Error(w, "Failed to create forwarded request", http.StatusInternalServerError)
		return
	}

	// Set the parameters
	forwardedRequest.URL.RawQuery = r.URL.RawQuery

	// Copy headers from the original request
	forwardedRequest.Header = make(http.Header)
	for key, values := range r.Header {
		for _, value := range values {
			forwardedRequest.Header.Add(key, value)
		}
	}

	// Make the request to the forward host
	client := http.Client{}
	resp, err := client.Do(forwardedRequest)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	log.Log.Debugf("Response of proxied request is %d\n", resp.StatusCode)

	// Write the response status code
	w.WriteHeader(resp.StatusCode)

	// Write the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
