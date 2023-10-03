package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLoginUsers(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_Seconds int    `json:"expires_in_seconds,omitempty"`
	}

	type responseBody struct {
		Id            int    `json:"id"`
		Email         string `json:"email"`
		Is_chirpy_red bool   `json:"is_chirpy_red"`
		Token         string `json:"token"`
		Refresh_Token string `json:"refresh_token"`
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
		return
	}

	for _, user := range usersDB {
		if req.Email == user.Email {

			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
			if err != nil {
				respondWithJSON(w, 401, "Invalid password")
				return
			}

			id := user.Id
			loginToken, err := cfg.createJWT(id, false)
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Unable to generate JWT token, %s", err))
				return
			}

			refreshToken, err := cfg.createJWT(id, true)
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Unable to generate JWT token, %s", err))
				return
			}

			res := responseBody{
				Id:            user.Id,
				Email:         user.Email,
				Is_chirpy_red: user.Is_chirpy_red,
				Token:         loginToken,
				Refresh_Token: refreshToken,
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
