package main

import (
	"log"

	"github.com/brenombrezende/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	dbHandler := database.DB{}
	//dbHandler.EnsureDB()
	testChirp, err := dbHandler.LoadDB()
	if err != nil {
		log.Fatal(err)
	}
	err = dbHandler.WriteDB(testChirp)

	if err != nil {
		log.Fatal(err)
	}
	dbHandler.GetChirps()
	//serverStarter()
}
