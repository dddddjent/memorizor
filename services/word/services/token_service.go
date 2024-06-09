package services

import (
	"crypto/rsa"
	"log"
	"memorizor/services/word/model"
	"memorizor/services/word/util"
)

type sTokenService struct {
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	accessTokenTimeout  int64
}
type STokenServiceConfig struct {
	PrivateKey          *rsa.PrivateKey
	PublicKey           *rsa.PublicKey
	AccessTokenTimeout  int64
}

func NewSTokenService(config *STokenServiceConfig) ITokenService {
	return &sTokenService{
		privateKey:          config.PrivateKey,
		publicKey:           config.PublicKey,
		accessTokenTimeout:  config.AccessTokenTimeout,
	}
}

func (s *sTokenService) ValidateAccessToken(tokenString string) (*model.User, error) {
	user, err := util.ValidateAccessToken(tokenString, s.publicKey)
    log.Println("Valid access token")
	if err != nil {
		return nil, err
	}
	return user, nil
}
