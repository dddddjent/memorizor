package services

import "github.com/gofrs/uuid"
import "memorizor/services/account/model"

type IUserService interface {
	GetByUUID(uuid.UUID) (*model.User, error)
	SignUp(*model.User) error
	SignIn(*model.User) error
}

type ITokenService interface {
	CreatePairFromUser(user *model.User, prevToken uuid.UUID) (*model.TokenPair, error)
	ValidateAccessToken(tokenString string) (*model.User, error)
	ValidateRefreshToken(tokenString string) (*model.SRefreshToken, error)
}
