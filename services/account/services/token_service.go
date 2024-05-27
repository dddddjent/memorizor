package services

import (
	"crypto/rsa"
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/util"
	"time"

	"github.com/gofrs/uuid"
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

func (s *sTokenService) CreatePairFromUser(user *model.User, prevToken uuid.UUID) (*model.TokenPair, error) {
	if prevToken != uuid.Nil {
		if err := s.tokenRepository.DeleteRefreshToken(user.UUID, prevToken); err != nil {
			if err, ok := err.(*util.Error); ok {
				return nil, err
			}
			return nil, util.NewInternal("Unable to delete previous token.\n" + err.Error())
		}
	}

	accessTimeoutDuration := time.Duration(s.accessTokenTimeout) * time.Second
	refreshTimeoutDuration := time.Duration(s.refreshTokenTimeout) * time.Second

	accessToken, err := util.GenerateAccessToken(user, s.privateKey, accessTimeoutDuration)
	if err != nil {
		return nil, util.NewInternal("Could not generate access token")
	}

	refreshToken, err := util.GenerateRefreshToken(user.UUID, s.refreshSecret, refreshTimeoutDuration)
	if err != nil {
		return nil, util.NewInternal("Could not generate refresh token")
	}

	if err := s.tokenRepository.SetRefreshToken(user.UUID, refreshToken.ID, refreshTimeoutDuration); err != nil {
		if err, ok := err.(*util.Error); ok {
			return nil, err
		}
		return nil, util.NewInternal("Unable to set refresh token.\n" + err.Error())
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (s *sTokenService) ValidateAccessToken(tokenString string) (*model.User, error) {
	user, err := util.ValidateAccessToken(tokenString, s.publicKey)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *sTokenService) ValidateRefreshToken(tokenString string) (*model.SRefreshToken, error) {
	refreshToken, err := util.ValidateRefreshToken(tokenString, s.refreshSecret)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}
