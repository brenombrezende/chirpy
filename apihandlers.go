package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateApi(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Body string `json:"body,omitempty"`
	}

	type responseBody struct {
		Valid bool `json:"valid"`
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

	err = respondWithJSON(w, 200, responseBody{
		Valid: true,
	})
	if err != nil {
		log.Printf("Error - %v", err)
	}
}
