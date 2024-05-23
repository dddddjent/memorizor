package services

import "github.com/gofrs/uuid"
import "memorizor/services/account/model"

type IUserService interface {
	GetByUUID(uuid.UUID) (*model.User, error)
}
