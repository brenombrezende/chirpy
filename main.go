package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/brenombrezende/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	router := chi.NewRouter()
	routerApi := chi.NewRouter()
	routerAdmin := chi.NewRouter()

	router.Mount("/api", routerApi)
	router.Mount("/admin", routerAdmin)

	fsHandler := apiCfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)
	router.Handle("/", fsHandler)

	routerApi.Get("/healthz", handlerHealthCheck)
	routerApi.Get("/reset", apiCfg.handlerResetMetrics)
	routerApi.Get("/chirps/{chirpID}", apiCfg.handlerGetChirps)
	routerApi.Post("/chirps", apiCfg.handlerValidateChirp)
	routerApi.Post("/users", apiCfg.handlerValidateUsers)

	routerAdmin.HandleFunc("/metrics", apiCfg.handlerDisplayMetrics)

	corsMux := middlewareCors(router)
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
