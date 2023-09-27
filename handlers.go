package main

import (
	"fmt"
	"net/http"
)

func handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
	fmt.Printf("handlerReadiness endpoint called\n")
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
	cfg.DB.ClearChirps()
	w.Write([]byte("Metrics reset to 0"))
	fmt.Printf("handlerResetMetrics endpoint called - Current Hits = %v\n", cfg.fileserverHits)

}

func (cfg *apiConfig) handlerDisplayMetrics(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>"
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(htmlTemplate, cfg.fileserverHits)))
	fmt.Printf("handlerDisplayMetrics endpoint called - Current Hits = %v\n", cfg.fileserverHits)
}
