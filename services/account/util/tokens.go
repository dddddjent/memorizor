package util

import (
	"crypto/rsa"
	"memorizor/services/account/model"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateIDToken(user *model.User, key *rsa.PrivateKey) (string, error) {
	unixTime := time.Now().Unix()
	expireTime := unixTime + 60*15 // 15 mins
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

type SRefreshToken struct {
	TokenString string
	ID          string
	ExpiresIn   time.Duration
}

func GenerateRefreshToken(id uuid.UUID, secret string) (*SRefreshToken, error) {
	currentTime := time.Now()
	tokenExpireTime := currentTime.AddDate(0, 0, 3)
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
