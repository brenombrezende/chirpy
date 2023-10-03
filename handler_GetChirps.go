package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/brenombrezende/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	s := r.URL.Query().Get("author_id")

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
	}

	if s != "" {
		qId, _ := strconv.Atoi(s)
		filteredChirps := []database.Chirp{}
		for i := range chirps {
			if chirps[i].Author_Id == qId {
				filteredChirps = append(filteredChirps, chirps[i])
			}

		}
		chirps = filteredChirps
	}

	sorting := r.URL.Query().Get("sort")
	if sorting == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		})

	} else {

		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})

	}

	respondWithJSON(w, 200, chirps)
}
