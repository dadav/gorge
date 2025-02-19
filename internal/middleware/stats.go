package middleware

import (
	"net/http"
	"sync"
	"time"
)

type Statistics struct {
	ActiveConnections             int
	TotalConnections              int
	TotalResponseTime             time.Duration
	TotalCacheHits                int
	TotalCacheMisses              int
	ConnectionsPerEndpoint        map[string]int
	ResponseTimePerEndpoint       map[string]time.Duration
	CacheHitsPerEndpoint          map[string]int
	CacheMissesPerEndpoint        map[string]int
	Mutex                         sync.Mutex
	ProxiedConnections            int
	ProxiedConnectionsPerEndpoint map[string]int
}

func NewStatistics() *Statistics {
	return &Statistics{
		ActiveConnections:             0,
		TotalConnections:              0,
		TotalResponseTime:             0,
		TotalCacheHits:                0,
		TotalCacheMisses:              0,
		ConnectionsPerEndpoint:        make(map[string]int),
		CacheHitsPerEndpoint:          make(map[string]int),
		CacheMissesPerEndpoint:        make(map[string]int),
		ResponseTimePerEndpoint:       make(map[string]time.Duration),
		ProxiedConnections:            0,
		ProxiedConnectionsPerEndpoint: make(map[string]int),
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
				stats.Mutex.Unlock()
			}()

			next.ServeHTTP(w, r)
		})
	}
}
