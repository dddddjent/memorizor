package repository

import (
	"memorizor/services/account/model"
	"time"

	"github.com/gofrs/uuid"
)

type IUserRepository interface {
	FindByUUID(uuid.UUID) (*model.User, error)
	FindByUserName(string) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Create(*model.User) error
	// Update the user, and put the result in the input user
	Update(id uuid.UUID, update_map map[string]any) (*model.User, error)
}

type ITokenRepository interface {
	SetRefreshToken(userID, tokenID uuid.UUID, expiresIn time.Duration) error
	DeleteRefreshToken(userID, previousTokenID uuid.UUID) error
	DeleteUserRefreshTokens(userID uuid.UUID) error
}
