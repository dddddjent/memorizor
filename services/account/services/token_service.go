package services

import (
	"crypto/rsa"
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/util"
)

type sTokenService struct {
	tokenRepository     repository.ITokenRepository
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	refreshSecret       string
	idTokenTimeout      int64
	refreshTokenTimeout int64
}
type STokenServiceConfig struct {
	TokenRepository     repository.ITokenRepository
	PrivateKey          *rsa.PrivateKey
	PublicKey           *rsa.PublicKey
	RefreshSecret       string
	IdTokenTimeout      int64
	RefreshTokenTimeout int64
}

func NewSTokenService(config *STokenServiceConfig) ITokenService {
	return &sTokenService{
		tokenRepository:     config.TokenRepository,
		privateKey:          config.PrivateKey,
		publicKey:           config.PublicKey,
		refreshSecret:       config.RefreshSecret,
		idTokenTimeout:      config.IdTokenTimeout,
		refreshTokenTimeout: config.RefreshTokenTimeout,
	}
}

func (s *sTokenService) CreatePairFromUser(user *model.User, prevToken string) (*model.TokenPair, error) {
	idToken, err := util.GenerateIDToken(user, s.privateKey, s.idTokenTimeout)
	if err != nil {
		return nil, &util.Error{Type: util.Internal, Message: "Could not generate id token"}
	}

	refreshToken, err := util.GenerateRefreshToken(user.UUID, s.refreshSecret, s.refreshTokenTimeout)
	if err != nil {
		return nil, &util.Error{Type: util.Internal, Message: "Could not generate refresh token"}
	}

	if err := s.tokenRepository.SetRefreshToken(user.UUID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		if err, ok := err.(*util.Error); ok {
			return nil, err
		}
		return nil, &util.Error{
			Type:    util.Internal,
			Message: "Unable to set refresh token.\n" + err.Error(),
		}
	}
	if prevToken != "" {
		if err := s.tokenRepository.DeleteRefreshToken(user.UUID.String(), prevToken); err != nil {
			if err, ok := err.(*util.Error); ok {
				return nil, err
			}
			return nil, &util.Error{
				Type:    util.Internal,
				Message: "Unable to delete previous token.\n" + err.Error(),
			}
		}
	}

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.TokenString,
	}, nil
}
