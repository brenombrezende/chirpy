package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{}
	r := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	r.Get("/healthz", handlerReadiness)
	r.Get("/reset", apiCfg.handlerResetMetrics)
	r.HandleFunc("/metrics", apiCfg.handlerDisplayMetrics)
	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	log.Printf("Starting server on path %s and port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
