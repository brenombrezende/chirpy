package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id        int    `json:"id"`
	Body      string `json:"body"`
	Author_Id int    `json:"author_id"`
}

type User struct {
	Id            int    `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Is_chirpy_red bool   `json:"is_chirpy_red"`
}

type RevokedToken struct {
	Id        string    `json:"id"`
	RevokedAt time.Time `json:"revoked_at"`
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RevokedTokens map[string]RevokedToken `json:"RevokedTokens"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RevokedTokens: map[string]RevokedToken{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbS := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if err != nil {
		return dbS, err
	}
	json.Unmarshal(dat, &dbS)
	return dbS, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) deleteDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	err := os.Remove(db.path)
	if err != nil {
		return err
	}

	return nil
}
