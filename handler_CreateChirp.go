package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Body string `json:"body,omitempty"`
	}

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	token, err := cfg.validateJWT(reqToken)
	if err != nil {
		respondWithError(w, 401, "Internal server error")
		return
	}

	author_id := token.(TokenInfoStruct).id

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err = decoder.Decode(&req)
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

	chirp, err := cfg.DB.CreateChirp(clearedString, author_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to write chirp into disk")
	}

	err = respondWithJSON(w, 201, chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
	}
}
