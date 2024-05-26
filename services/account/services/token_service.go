package services

import (
	"crypto/rsa"
	"log"
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/util"
)

type sTokenService struct {
	tokenRepository     repository.ITokenRepository
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	refreshSecret       string
	accessTokenTimeout  int64
	refreshTokenTimeout int64
}
type STokenServiceConfig struct {
	TokenRepository     repository.ITokenRepository
	PrivateKey          *rsa.PrivateKey
	PublicKey           *rsa.PublicKey
	RefreshSecret       string
	AccessTokenTimeout  int64
	RefreshTokenTimeout int64
}

func NewSTokenService(config *STokenServiceConfig) ITokenService {
	return &sTokenService{
		tokenRepository:     config.TokenRepository,
		privateKey:          config.PrivateKey,
		publicKey:           config.PublicKey,
		refreshSecret:       config.RefreshSecret,
		accessTokenTimeout:  config.AccessTokenTimeout,
		refreshTokenTimeout: config.RefreshTokenTimeout,
	}
}

func (s *sTokenService) CreatePairFromUser(user *model.User, prevToken string) (*model.TokenPair, error) {
	accessToken, err := util.GenerateAccessToken(user, s.privateKey, s.accessTokenTimeout)
	if err != nil {
		return nil, util.NewInternal("Could not generate access token")
	}

	refreshToken, err := util.GenerateRefreshToken(user.UUID, s.refreshSecret, s.refreshTokenTimeout)
	if err != nil {
		return nil, util.NewInternal("Could not generate refresh token")
	}

	if err := s.tokenRepository.SetRefreshToken(user.UUID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		if err, ok := err.(*util.Error); ok {
			return nil, err
		}
		return nil, util.NewInternal("Unable to set refresh token.\n" + err.Error())
	}
	if prevToken != "" {
		if err := s.tokenRepository.DeleteRefreshToken(user.UUID.String(), prevToken); err != nil {
			if err, ok := err.(*util.Error); ok {
				return nil, err
			}
			return nil, util.NewInternal("Unable to delete previous token.\n" + err.Error())
		}
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.TokenString,
	}, nil
}

func (s *sTokenService) ValidateAccessToken(tokenString string) (*model.User, error) {
	claims, err := util.ValidateAccessToken(tokenString, s.publicKey)
	if err != nil {
		log.Println(err.Error())
		log.Println("Could not get claims from the request")
		return nil, util.NewAuthorization("Unable to verify user from the access token")
	}
	log.Println("claims: ", claims)
	user := claims.User
	return user, nil
}
