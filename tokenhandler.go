package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Issuer    string
	IssuedAt  jwt.NumericDate
	ExpiresAt jwt.NumericDate
	Subject   string
}

func (cfg apiConfig) createJWT(id, expiration int) (token string, err error) {

	key := []byte(cfg.jwtSecret)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "Chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiration))),
			Subject:   strconv.Itoa(id),
		})
	token, err = t.SignedString(key)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (cfg apiConfig) validateJWT(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid")
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	convertedId, _ := strconv.Atoi(id)
	return convertedId, nil
}
