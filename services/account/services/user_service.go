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

func (service *sUserService) SignIn(user *model.User) error {
	password := user.Password
	userName := user.UserName
	email := user.Email
	if userName != "" && email != "" || userName == "" && email == "" {
		return util.NewBadRequest("Bad email and user name combo for sign in")
	}

	var (
		userFound *model.User
		err       error
	)
	if email != "" {
		userFound, err = service.repository.FindByEmail(email)
	} else {
		userFound, err = service.repository.FindByUserName(userName)
	}
	if err != nil {
		return util.NewAuthorization("No user found")
	}

	compareResult, err := util.ComparePassword(userFound.Password, password)
	if err != nil {
		return util.NewInternal("Password compare failed")
	}
	if compareResult == false {
		return util.NewAuthorization("Incorrect password")
	}
	userFound.DeepCopyTo(user)
	return nil
}
