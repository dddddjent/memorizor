package util

import (
	"crypto/rsa"
	"fmt"
	"memorizor/services/account/model"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	User *model.User `json:"user"`
}

func GenerateAccessToken(user *model.User, key *rsa.PrivateKey, timeout int64) (string, error) {
	unixTime := time.Now()
	expireTime := unixTime.Add(time.Duration(timeout) * time.Second) // 15 mins
	claims := AccessTokenClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: unixTime},
			ExpiresAt: &jwt.NumericDate{Time: expireTime},
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateAccessToken(tokenString string, key *rsa.PublicKey) (*AccessTokenClaims, error) {
	parsedClaims := AccessTokenClaims{}
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
type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UUID    string `json:"uuid"`
	TokenId string `json:"token_id"`
}

func GenerateRefreshToken(id uuid.UUID, secret string, timeout int64) (*SRefreshToken, error) {
	currentTime := time.Now()
	tokenExpireTime := currentTime.Add(time.Duration(timeout) * time.Second)
	tokenID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
    
	claims := RefreshTokenClaims{
		UUID:    id.String(),
		TokenId: tokenID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: currentTime},
			ExpiresAt: &jwt.NumericDate{Time: tokenExpireTime},
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
