package middleware

import (
	"net/http"
	"net/http/httputil"
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
			capturedResponseWriter := &capturedResponseWriter{ResponseWriter: w}
			next.ServeHTTP(capturedResponseWriter, r)

			if forwardToProxy(capturedResponseWriter.status) {
				log.Log.Infof("Forwarding request to %s\n", upstreamHost)
				u, err := url.Parse(upstreamHost)
				if err != nil {
					log.Log.Error(err)
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
					capturedResponseWriter.sendCapturedResponse()
				}

				proxy.ServeHTTP(w, r)
				return
			}

			// If the response status is not 404, serve the original response
			capturedResponseWriter.sendCapturedResponse()
		})
	}
}
