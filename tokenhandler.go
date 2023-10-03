package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenInfoStruct struct {
	id     int
	issuer string
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Issuer    string
	IssuedAt  jwt.NumericDate
	ExpiresAt jwt.NumericDate
	Subject   string
}

func (cfg apiConfig) createJWT(id int, isRefresh bool) (token string, err error) {
	issuerString := "chirpy-access"
	expirationTime := time.Now().Add(time.Hour * 1)

	if isRefresh {
		issuerString = "chirpy-refresh"
		expirationTime = time.Now().Add(time.Hour * 24 * 60)
	}

	key := []byte(cfg.jwtSecret)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuerString,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   strconv.Itoa(id),
		})
	token, err = t.SignedString(key)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (cfg apiConfig) validateJWT(tokenStr string) (interface{}, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid")
	}

	tokeninfo := TokenInfoStruct{}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	tokeninfo.id, _ = strconv.Atoi(id)

	tokeninfo.issuer, err = token.Claims.GetIssuer()
	if err != nil {
		return 0, err
	}

	return tokeninfo, nil
}
