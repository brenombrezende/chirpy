package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerTokenRevoker(w http.ResponseWriter, r *http.Request) {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	token, err := cfg.validateJWT(reqToken)
	if err != nil {
		respondWithError(w, 401, "Internal server error")
		return
	}

	if token.(TokenInfoStruct).issuer != "chirpy-refresh" {
		respondWithJSON(w, 401, "Unauthorized")
		return
	}

	err = cfg.DB.RevokeToken(reqToken)
	if err != nil {
		respondWithError(w, 501, "Internal server error")
		return
	}

	err = respondWithJSON(w, 200, "")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return

	}
}
