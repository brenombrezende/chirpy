package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Body string `json:"body,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithJSON(w, 500, requestBody{
			Body: "Something went wrong",
		})

		log.Printf("Error decoding parameters: %s", err)
		return
	}

	if len(req.Body) > 140 {
		respondWithJSON(w, 400, requestBody{
			Body: "Chirp is too long",
		})

		return
	}

	clearedString := profanityChecker(string(req.Body))

	chirp, err := cfg.DB.CreateChirp(clearedString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to write chirp into disk")
	}

	err = respondWithJSON(w, 201, chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
	}
}
