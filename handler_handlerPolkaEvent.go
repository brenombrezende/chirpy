package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerPolkaEvent(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		respondWithJSON(w, 401, "")
		return
	}
	splitToken := strings.Split(reqToken, "ApiKey ")
	reqToken = splitToken[1]

	if reqToken != cfg.polkaKey {
		respondWithJSON(w, 401, "")
		return
	}

	errNotFound := errors.New("user not found")

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode request")
		return
	}

	if req.Event != "user.upgraded" {
		respondWithJSON(w, 200, "")
		return
	}

	err = cfg.DB.UpgradeUser(req.Data.UserID)
	if errors.Is(err, errNotFound) {
		respondWithError(w, 404, "")
		return
	}
	respondWithJSON(w, 200, "")

}
