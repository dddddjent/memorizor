package services

import (
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/util"

	"github.com/gofrs/uuid"
)

type sUserService struct {
	repository repository.IUserRepository
}

type SUserServiceConfig struct {
	Repository repository.IUserRepository
}

func NewSUserService(config *SUserServiceConfig) IUserService {
	return &sUserService{
		repository: config.Repository,
	}
}

func (service *sUserService) GetByUUID(id uuid.UUID) (*model.User, error) {
	return service.repository.FindByUUID(id)
}
func (service *sUserService) SignUp(user *model.User) error {
	encoded, err := util.EncodePassword(user.Password)
	if err != nil {
		return util.NewInternal("Failed to encode password\n" + err.Error())
	}

	user.Password = encoded
	if err := service.repository.Create(user); err != nil {
		return err
	}
	return nil
}
