package database

import "errors"

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

func (db *DB) ClearChirps() error {
	err := db.deleteDB()
	if err != nil {
		return err
	}
	db.ensureDB()
	return nil
}
