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
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Id:        id,
		Body:      body,
		Author_Id: author_id,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chirpID, author_id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	if dbStructure.Chirps[chirpID].Author_Id != author_id {
		return errors.New("not authorized")
	}

	delete(dbStructure.Chirps, chirpID)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}
	for i := range data.Chirps {
		chirps = append(chirps, data.Chirps[i])
	}
	return chirps, nil
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if email == user.Email {
			return User{}, errors.New("this email already exists")
		}
	}

	id := len(dbStructure.Users) + 1
	user := User{
		Id:       id,
		Email:    email,
		Password: password,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) RevokeToken(tokenstring string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	currentDate := time.Now()

	token := RevokedToken{
		Id:        tokenstring,
		RevokedAt: currentDate,
	}
	dbStructure.RevokedTokens[tokenstring] = token

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRevokedTokens() ([]RevokedToken, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	allTokens := []RevokedToken{}
	for i := range data.RevokedTokens {
		allTokens = append(allTokens, data.RevokedTokens[i])
	}
	return allTokens, nil
}

func (db *DB) UpdateUser(email, password string, id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	updatedUser := User{
		Email:    email,
		Password: password,
	}

	for _, user := range dbStructure.Users {
		if id == user.Id {
			dbStructure.Users[user.Id] = updatedUser
		}
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

// GetChirps returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := []User{}
	for i := range data.Users {
		users = append(users, data.Users[i])
	}
	return users, nil
}

func (db *DB) ClearChirps() error {
	err := db.deleteDB()
	if err != nil {
		return err
	}
	db.ensureDB()
	return nil
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
