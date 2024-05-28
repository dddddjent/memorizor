package services

import "github.com/gofrs/uuid"
import "memorizor/services/account/model"

type IUserService interface {
	GetByUUID(uuid.UUID) (*model.User, error)
	SignUp(*model.User) error
	SignIn(*model.User) error
	Update(id uuid.UUID, update_map map[string]any) (*model.User, error)
}

type ITokenService interface {
	CreatePairFromUser(user *model.User, prevToken uuid.UUID) (*model.TokenPair, error)
	ValidateAccessToken(tokenString string) (*model.User, error)
	ValidateRefreshToken(tokenString string) (*model.SRefreshToken, error)
	// If prevToken is not Nil, then only delete that refresh token
	SignOut(user *model.User, prevToken uuid.UUID) error
}
