package main

import (
	"fmt"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
	fmt.Printf("handlerReadiness endpoint called\n")
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
