package util

import (
	"crypto/rsa"
	"log"
	"memorizor/services/account/model"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type accessTokenClaims struct {
	jwt.RegisteredClaims
	User *model.User `json:"user"`
}

func GenerateAccessToken(user *model.User, key *rsa.PrivateKey, timeout time.Duration) (string, error) {
	currentTime := time.Now()
	expireTime := currentTime.Add(timeout) // 15 mins
	claims := accessTokenClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: currentTime},
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

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	UUID    string `json:"uuid"`
	TokenId string `json:"token_id"`
}

func GenerateRefreshToken(id uuid.UUID, secret string, timeout time.Duration) (*model.SRefreshToken, error) {
	currentTime := time.Now()
	tokenExpireTime := currentTime.Add(timeout)
	tokenID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	claims := refreshTokenClaims{
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
	return &model.SRefreshToken{
		TokenString: tokenString,
		ID:          tokenID,
		UUID:        id,
	}, nil
}

func ValidateRefreshToken(tokenString string, secret string) (*model.SRefreshToken, error) {
	parsedClaims := refreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &parsedClaims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Println(err.Error())
		log.Println("Could not get claims from the request")
		return nil, NewAuthorization("Unable to verify user from the refresh token: " + err.Error())
	}
	if !token.Valid {
		return nil, NewAuthorization("Refresh token is invalid")
	}
	log.Println("claims: ", parsedClaims)
	tokenID, _ := uuid.FromString(parsedClaims.TokenId)
	userID, _ := uuid.FromString(parsedClaims.UUID)
	refreshToken := model.SRefreshToken{
		TokenString: tokenString,
		ID:          tokenID,
		UUID:        userID,
	}
	return &refreshToken, nil
}
