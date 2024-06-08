package util

import (
	"crypto/rsa"
	"log"
	"memorizor/services/word/model"

	"github.com/golang-jwt/jwt/v5"
)

type accessTokenClaims struct {
	jwt.RegisteredClaims
	User *model.User `json:"user"`
}

func ValidateAccessToken(tokenString string, key *rsa.PublicKey) (*model.User, error) {
	parsedClaims := accessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &parsedClaims, func(t *jwt.Token) (any, error) {
		return key, nil
	})
	if err != nil {
		log.Println(err.Error())
		log.Println("Could not get claims from the request")
		return nil, NewAuthorization("Unable to verify user from the access token: " + err.Error())
	}
	if !token.Valid {
		return nil, NewAuthorization("Access token is invalid")
	}
	log.Println("claims: ", parsedClaims)
	return parsedClaims.User, nil
}
