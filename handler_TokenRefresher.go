package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerTokenRefresher(w http.ResponseWriter, r *http.Request) {

	type responseBody struct {
		Token string `json:"token"`
	}

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	allRevokedTokens, err := cfg.DB.GetRevokedTokens()
	if err != nil {
		respondWithError(w, 501, "Internal server error")
		return
	}

	for _, revokedToken := range allRevokedTokens {
		if reqToken == revokedToken.Id {
			respondWithError(w, 401, "Expired TOken")
			return
		}
	}

	token, err := cfg.validateJWT(reqToken)
	if err != nil {
		respondWithError(w, 401, "Internal server error")
		return
	}

	if token.(TokenInfoStruct).issuer != "chirpy-refresh" {
		respondWithJSON(w, 401, "Unauthorized")
		return
	}

	newToken, err := cfg.createJWT(token.(TokenInfoStruct).id, true)
	if err != nil {
		respondWithError(w, 401, "Internal server error")
		return
	}

	res := responseBody{
		Token: newToken,
	}

	err = respondWithJSON(w, 200, res)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return

	}
}
