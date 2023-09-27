package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
	}

	id := chi.URLParam(r, "chirpID")
	idInt, _ := strconv.Atoi(id)

	for i := range chirps {
		if chirps[i].Id == idInt {
			response := chirps[i]
			respondWithJSON(w, 200, response)
			return
		}
	}
	respondWithJSON(w, 404, "")
}
