package middleware

import (
	"net/http"
	"sync"
	"time"
)

type Statistics struct {
	ActiveConnections             int
	TotalConnections              int
	ProxiedConnections            int
	TotalResponseTime             time.Duration
	ConnectionsPerEndpoint        map[string]int
	ProxiedConnectionsPerEndpoint map[string]int
	ResponseTimePerEndpoint       map[string]time.Duration
	Mutex                         sync.Mutex
}

func NewStatistics() *Statistics {
	return &Statistics{
		ActiveConnections:             0,
		TotalConnections:              0,
		ProxiedConnections:            0,
		TotalResponseTime:             0,
		ConnectionsPerEndpoint:        make(map[string]int),
		ProxiedConnectionsPerEndpoint: make(map[string]int),
		ResponseTimePerEndpoint:       make(map[string]time.Duration),
	}
}

func StatisticsMiddleware(stats *Statistics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			stats.Mutex.Lock()
			stats.ActiveConnections++
			stats.TotalConnections++
			stats.ConnectionsPerEndpoint[r.URL.Path]++
			stats.Mutex.Unlock()

			defer func() {
				duration := time.Since(start)
				stats.Mutex.Lock()
				stats.ActiveConnections--
				stats.TotalResponseTime += duration
				stats.ResponseTimePerEndpoint[r.URL.Path] += duration

				if w.Header().Get("X-Proxied-To") != "" {
					stats.ProxiedConnections++
					stats.ProxiedConnectionsPerEndpoint[r.URL.Path]++
				}

				stats.Mutex.Unlock()
			}()

			next.ServeHTTP(w, r)
		})
	}
}
