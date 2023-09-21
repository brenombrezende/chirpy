package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

var dbName string = "database.json"

func CreateChirp(body string) (Chirp, error) {
	c := Chirp{}

	err := os.WriteFile(dbName, []byte(body), 0666)
	if err != nil {
		log.Fatal(err)
	}
	return c, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	test := []Chirp{}
	data, err := db.LoadDB()
	if err != nil {
		log.Printf("Error - unable to load chirps in GetChirps - %s", err)
		return nil, err
	}

	for i, _ := range data.Chirps {
		test = append(test, data.Chirps[i])
	}
	return test, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) EnsureDB() error {

	if _, err := os.Stat(dbName); err != nil {
		_, err := os.Create(dbName)
		if err != nil {
			log.Fatalf("Error creating file - %s", err)
			return err
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) LoadDB() (DBStructure, error) {
	dbS := DBStructure{}
	dat, err := os.ReadFile(dbName)
	if err != nil {
		log.Printf("Unable to load DB from file %s - Error %s", dbName, err)
		return dbS, err
	}
	json.Unmarshal(dat, &dbS)
	return dbS, nil
}

// writeDB writes the database file to disk
func (db *DB) WriteDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		log.Printf("ERROR - Unable to Marshal data to disk - %s", err)
		return err
	}
	err = os.WriteFile(dbName, data, 0666)

	if err != nil {
		log.Printf("ERROR - Unable to write data to disk - %s", err)
		return err
	}

	return nil
}
