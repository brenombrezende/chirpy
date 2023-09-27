package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerValidateUsers(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode request")
		return
	}

	user, err := cfg.DB.CreateUser(req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to write user into disk")
		return
	}

	err = respondWithJSON(w, 201, user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return

	}
}
