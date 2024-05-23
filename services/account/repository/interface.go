package repository

import (
	"memorizor/services/account/model"

	"github.com/gofrs/uuid"
)

type IUserRepository interface {
	FindByUUID(uuid.UUID) (*model.User, error)
}
