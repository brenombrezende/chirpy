package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email    string `json:"email"`
		Password string `jason:"password"`
	}

	type responseBody struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode request")
		return
	}

	usersDB, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch users")
	}

	for _, user := range usersDB {
		if req.Email == user.Email {

			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
			if err != nil {
				respondWithJSON(w, 401, "Invalid password")
				return
			}

			res := responseBody{
				Id:    user.Id,
				Email: user.Email,
			}

			err = respondWithJSON(w, 200, res)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
				return

			}
			return
		}
	}
	respondWithJSON(w, 403, "Not authorized.")
}
