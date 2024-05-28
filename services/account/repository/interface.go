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
}

type ITokenRepository interface {
	SetRefreshToken(userID, tokenID uuid.UUID, expiresIn time.Duration) error
	DeleteRefreshToken(userID, previousTokenID uuid.UUID) error
	DeleteUserRefreshTokens(userID uuid.UUID) error
}
