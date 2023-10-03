package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerPasswordChange(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email    string `json:"email"`
		Password string `jason:"password"`
	}

	type responseBody struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	token, err := cfg.validateJWT(reqToken)
	if err != nil {
		respondWithError(w, 401, "Internal server error")
		return
	}

	if token.(TokenInfoStruct).issuer == "chirpy-refresh" {
		respondWithJSON(w, 401, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode request")
		return
	}

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to properly hash password")
		return
	}

	user, err := cfg.DB.UpdateUser(req.Email, string(hashedPW), token.(TokenInfoStruct).id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to create user - %s", err))
		return
	}

	res := responseBody{
		Email: user.Email,
		Id:    token.(TokenInfoStruct).id,
	}

	err = respondWithJSON(w, 200, res)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return

	}
}
