package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		fmt.Printf("middlewareMetricsInc used - Current Hits = %v\n", cfg.fileserverHits)
		next.ServeHTTP(w, r)

	})
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
	w.Write([]byte("Metrics reset to 0"))
	fmt.Printf("handlerResetMetrics endpoint called - Current Hits = %v\n", cfg.fileserverHits)

}

func (cfg *apiConfig) handlerDisplayMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
	fmt.Printf("handlerDisplayMetrics endpoint called - Current Hits = %v\n", cfg.fileserverHits)
}
