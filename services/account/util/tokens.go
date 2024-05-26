package util

import (
	"crypto/rsa"
	"fmt"
	"memorizor/services/account/model"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user *model.User, key *rsa.PrivateKey, timeout int64) (string, error) {
	unixTime := time.Now().Unix()
	expireTime := unixTime + timeout // 15 mins
	t := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iat":  unixTime,
			"exp":  expireTime,
			"user": user,
		})
	tokenString, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateAccessToken(tokenString string, key *rsa.PublicKey) (*jwt.MapClaims, error) {
	parsedClaims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, &parsedClaims, func(t *jwt.Token) (any, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Access token is invalid")
	}
	return &parsedClaims, nil
}

type SRefreshToken struct {
	TokenString string
	ID          string
	ExpiresIn   time.Duration
}

func GenerateRefreshToken(id uuid.UUID, secret string, timeout int64) (*SRefreshToken, error) {
	currentTime := time.Now()
	tokenExpireTime := currentTime.Add(time.Duration(timeout) * time.Second)
	tokenID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iat":      currentTime.Unix(),
			"exp":      tokenExpireTime.Unix(),
			"token_id": tokenID.String(),
			"uuid":     id.String(),
		})
	tokenString, err := t.SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}
	return &SRefreshToken{
		TokenString: tokenString,
		ID:          tokenID.String(),
		ExpiresIn:   tokenExpireTime.Sub(currentTime),
	}, nil
}
