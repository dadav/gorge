package middleware

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dadav/gorge/internal/log"
)

// capturedResponseWriter is a custom response writer that captures the response status
type capturedResponseWriter struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func NewCapturedResponseWriter(w http.ResponseWriter) *capturedResponseWriter {
	return &capturedResponseWriter{
		ResponseWriter: w,
		body:           new(bytes.Buffer),
	}
}

func (w *capturedResponseWriter) WriteHeader(code int) {
	w.status = code
}

func (w *capturedResponseWriter) Write(body []byte) (int, error) {
	return w.body.Write(body)
}

func (w *capturedResponseWriter) sendCapturedResponse() {
	w.ResponseWriter.WriteHeader(w.status)
	w.ResponseWriter.Write(w.body.Bytes())
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
		},
	}
}

func ProxyFallback(upstreamHost string, forwardToProxy func(int) bool, proxiedResponseCb func(*http.Response)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Store original headers before any modifications
			originalHeaders := make(http.Header)
			for k, v := range w.Header() {
				originalHeaders[k] = v
			}

			capturedResponseWriter := NewCapturedResponseWriter(w)
			next.ServeHTTP(capturedResponseWriter, r)

			if forwardToProxy(capturedResponseWriter.status) {
				log.Log.Infof("Forwarding request to %s\n", upstreamHost)
				u, err := url.Parse(upstreamHost)
				if err != nil {
					log.Log.Error(err)
					// Restore original headers before sending captured response
					for k, v := range originalHeaders {
						w.Header()[k] = v
					}
					capturedResponseWriter.sendCapturedResponse()
					return
				}

				for k := range w.Header() {
					w.Header().Del(k)
				}

				proxy := NewSingleHostReverseProxy(u)

				proxy.ModifyResponse = func(r *http.Response) error {
					proxiedResponseCb(r)
					return nil
				}

				// if some error occurs, return the original content
				proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
					log.Log.Error(err)
					// Restore original headers before sending captured response
					for k, v := range originalHeaders {
						w.Header()[k] = v
					}
					capturedResponseWriter.sendCapturedResponse()
				}

				proxy.ServeHTTP(w, r)
				return
			}

			// If the response status is not 404, serve the original response
			// Restore original headers before sending captured response
			for k, v := range originalHeaders {
				w.Header()[k] = v
			}
			capturedResponseWriter.sendCapturedResponse()
		})
	}
}
