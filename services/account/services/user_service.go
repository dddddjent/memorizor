package services

import (
	"memorizor/services/account/model"
	"memorizor/services/account/repository"

	"github.com/gofrs/uuid"
)

type SUserService struct {
	repository repository.IUserRepository
}

type SUserServiceConfig struct {
	Repository repository.IUserRepository
}

func NewSUserService(config *SUserServiceConfig) IUserService {
	return &SUserService{
		repository: config.Repository,
	}
}

func (service *SUserService) GetByUUID(id uuid.UUID) (*model.User, error) {
	return service.repository.FindByUUID(id)
}
