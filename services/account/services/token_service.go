package services

import (
	"crypto/rsa"
	"memorizor/services/account/model"
	"memorizor/services/account/util"
)

type sTokenService struct {
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	refreshSecret       string
	idTokenTimeout      int64
	refreshTokenTimeout int64
}
type STokenServiceConfig struct {
	PrivateKey          *rsa.PrivateKey
	PublicKey           *rsa.PublicKey
	RefreshSecret       string
	IdTokenTimeout      int64
	RefreshTokenTimeout int64
}

func NewSTokenService(config *STokenServiceConfig) ITokenService {
	return &sTokenService{
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
		return nil, &util.Error{Type: util.Internal, Message: "Can't generate id token"}
	}

	refreshToken, err := util.GenerateRefreshToken(user.UUID, s.refreshSecret, s.refreshTokenTimeout)
	if err != nil {
		return nil, &util.Error{Type: util.Internal, Message: "Can't generate refresh token"}
	}

	// TODO: store refresh token by calling token repo's method

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.TokenString,
	}, nil
}
