package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dadav/gorge/internal/config"
	"github.com/dadav/gorge/internal/log"
)

type ContentHeaders struct {
	Type        string `json:"type"`
	Encoding    string `json:"encoding"`
	Disposition string `json:"disposition"`
}

func CacheMiddleware(prefixes []string, cacheDir string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(r.URL.Path, prefix) {
					matched = true
					break
				}
			}

			if !matched || r.Method != "GET" {
				next.ServeHTTP(w, r)
				return
			}

			cacheKeyRaw := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)
			hash := sha256.New()
			hash.Write([]byte(cacheKeyRaw))
			cacheKeyHash := hex.EncodeToString(hash.Sum(nil))
			cacheFilePath := filepath.Join(cacheDir, cacheKeyHash)
			cacheFileHeadersPath := fmt.Sprintf("%s_headers", cacheFilePath)

			cacheControlHeader := r.Header.Get("Cache-Control")
			if !strings.Contains(cacheControlHeader, "no-cache") {
				if cacheFileInfo, err := os.Stat(cacheFilePath); err == nil {
					expirationTime := cacheFileInfo.ModTime().Add(time.Duration(config.CacheMaxAge) * time.Second)
					if time.Now().After(expirationTime) {
						log.Log.Debugf("Cached file expired: %s\n", cacheFilePath)
						err := os.Remove(cacheFilePath)
						if err != nil {
							log.Log.Error(err)
						}
					} else {
						data, err := os.ReadFile(cacheFilePath)
						if err == nil {
							log.Log.Debugf("Send response from cache for %s\n", r.URL.Path)
							headerBytes, err := os.ReadFile(cacheFileHeadersPath)
							if err == nil {
								var contentHeaders ContentHeaders
								json.Unmarshal(headerBytes, &contentHeaders)
								if contentHeaders.Type != "" {
									w.Header().Add("Content-Type", contentHeaders.Type)
								}
								if contentHeaders.Encoding != "" {
									w.Header().Add("Content-Encoding", contentHeaders.Encoding)
								}
								if contentHeaders.Disposition != "" {
									w.Header().Add("Content-Disposition", contentHeaders.Disposition)
								}
							}
							w.Write(data)
							return
						}

					}
				}
			}

			capturedResponseWriter := &capturedResponseWriter{ResponseWriter: w}
			next.ServeHTTP(capturedResponseWriter, r)

			if capturedResponseWriter.status == http.StatusOK && !strings.Contains(cacheControlHeader, "no-store") {
				err := os.WriteFile(cacheFilePath, capturedResponseWriter.body, 0600)
				if err != nil {
					log.Log.Error(err)
				}

				contentHeaders := ContentHeaders{
					Type:        capturedResponseWriter.Header().Get("Content-Type"),
					Encoding:    capturedResponseWriter.Header().Get("Content-Encoding"),
					Disposition: capturedResponseWriter.Header().Get("Content-Disposition"),
				}
				contentHeadersBytes, err := json.Marshal(contentHeaders)
				if err == nil {
					err = os.WriteFile(cacheFileHeadersPath, contentHeadersBytes, 0600)
					if err != nil {
						log.Log.Error(err)
					}
				}
			}

			capturedResponseWriter.sendCapturedResponse()
		})
	}
}
