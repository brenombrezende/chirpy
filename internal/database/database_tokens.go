package database

import "time"

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
