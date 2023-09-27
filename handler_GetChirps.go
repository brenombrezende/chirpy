package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		log.Print("Error getting chirps - %s", err)
	}

	json.Marshal(chirps)

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })

	respondWithJSON(w, 200, chirps)

}
