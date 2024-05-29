package services

import (
	"memorizor/services/account/model"
	"mime/multipart"

	"github.com/gofrs/uuid"
)

type IUserService interface {
	GetByUUID(uuid.UUID) (*model.User, error)
	SignUp(*model.User) error
	SignIn(*model.User) error
	Update(id uuid.UUID, update_map map[string]any) (*model.User, error)
	UpdateProfileImage(id uuid.UUID, imageFile multipart.File, imageType string) (imageURL string, err error)
}

type ITokenService interface {
	CreatePairFromUser(user *model.User, prevToken uuid.UUID) (*model.TokenPair, error)
	ValidateAccessToken(tokenString string) (*model.User, error)
	ValidateRefreshToken(tokenString string) (*model.SRefreshToken, error)
	// If prevToken is not Nil, then only delete that refresh token
	SignOut(user *model.User, prevToken uuid.UUID) error
}
