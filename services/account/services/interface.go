package services

import "github.com/gofrs/uuid"
import "memorizor/services/account/model"

type UserService interface {
	GetByUUID(uid uuid.UUID) (*model.User, error)
}
