package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	token, err := cfg.validateJWT(reqToken)
	if err != nil {
		respondWithError(w, 401, "Internal server error")
		return
	}

	author_id := token.(TokenInfoStruct).id

	id := chi.URLParam(r, "chirpID")
	idInt, _ := strconv.Atoi(id)

	err = cfg.DB.DeleteChirp(idInt, author_id)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Error - %s", err))
	}

	respondWithJSON(w, 200, "")
}
