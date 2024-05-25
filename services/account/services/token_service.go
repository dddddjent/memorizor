package services

import (
	"crypto/rsa"
	"memorizor/services/account/model"
	"memorizor/services/account/util"
)

type STokenService struct {
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	refreshSecret string
}
type STokenServiceConfig struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

func NewSTokenService(config *STokenServiceConfig) *STokenService {
	return &STokenService{
		privateKey:    config.PrivateKey,
		publicKey:     config.PublicKey,
		refreshSecret: config.RefreshSecret,
	}
}

func (s *STokenService) CreatePairFromUser(user *model.User, prevToken string) (*model.TokenPair, error) {
	idToken, err := util.GenerateIDToken(user, s.privateKey)
	if err != nil {
		return nil, &util.Error{Type: util.Internal, Message: "Can't generate id token"}
	}

	refreshToken, err := util.GenerateRefreshToken(user.UUID, s.refreshSecret)
	if err != nil {
		return nil, &util.Error{Type: util.Internal, Message: "Can't generate refresh token"}
	}

	// TODO: store refresh token by calling token repo's method

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.TokenString,
	}, nil
}
