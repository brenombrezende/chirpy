package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handlerValidateApi(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("handlerValidateApi called\n")

	type params struct {
		Body string `json:"body,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}
	err := decoder.Decode(&p)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		dat := params{
			Body: "Something went wrong",
		}
		resp, _ := json.Marshal(dat)
		w.Write(resp)

		log.Printf("Error decoding parameters: %s", err)
		return
	}

	if len(p.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		dat := params{
			Body: "Chirp is too long",
		}
		resp, _ := json.Marshal(dat)
		w.Write(resp)

		log.Printf("Body too long - %v characters: ", len(p.Body))
		return
	}

	validated := struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	}

	dat, _ := json.Marshal(validated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
