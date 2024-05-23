package services

import (
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/util"

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
func (service *SUserService) SignUp(user *model.User) error {
	encoded, err := util.EncodePassword(user.Password)
	if err != nil {
		return &util.Error{
			Type:    util.Internal,
			Message: "Failed to encode password",
		}
	}

	user.Password = encoded
	id, _ := uuid.NewV7()
	user.UUID = id
	if err := service.repository.Create(user); err != nil {
		return err
	}
	return nil
}
