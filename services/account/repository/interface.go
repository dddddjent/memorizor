package repository

import (
	"memorizor/services/account/model"
	"time"

	"github.com/gofrs/uuid"
)

type IUserRepository interface {
	FindByUUID(uuid.UUID) (*model.User, error)
	Create(*model.User) error
}

type ITokenRepository interface {
	SetRefreshToken(userID, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(userID, previousTokenID string) error
}
