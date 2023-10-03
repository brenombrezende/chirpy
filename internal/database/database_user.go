package database

import "errors"

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
		Id:            id,
		Email:         email,
		Password:      password,
		Is_chirpy_red: false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(email, password string, id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if id == user.Id {
			dbStructure.Users[id] = User{
				Email:         email,
				Password:      password,
				Is_chirpy_red: dbStructure.Users[id].Is_chirpy_red,
			}
		}
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return dbStructure.Users[id], nil
}

func (db *DB) UpgradeUser(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")
	}

	user.Is_chirpy_red = true

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
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
