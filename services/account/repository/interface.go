package repository

import (
	"memorizor/services/account/model"

	"github.com/gofrs/uuid"
)

type UserRepository interface {
	FindByUUID(uid uuid.UUID) (*model.User, error)
}
