package services

import (
	"fmt"
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/util"
	"mime/multipart"

	"github.com/gofrs/uuid"
)

type sUserService struct {
	userRepo         repository.IUserRepository
	profileImageRepo repository.IProfileImageRepository
}

type SUserServiceConfig struct {
	UserRepository   repository.IUserRepository
	ProfileImageRepo repository.IProfileImageRepository
}

func NewSUserService(config *SUserServiceConfig) IUserService {
	return &sUserService{
		userRepo:         config.UserRepository,
		profileImageRepo: config.ProfileImageRepo,
	}
}

func (service *sUserService) GetByUUID(id uuid.UUID) (*model.User, error) {
	return service.userRepo.FindByUUID(id)
}
func (service *sUserService) SignUp(user *model.User) error {
	encoded, err := util.EncodePassword(user.Password)
	if err != nil {
		return util.NewInternal("Failed to encode password\n" + err.Error())
	}

	user.Password = encoded
	if err := service.userRepo.Create(user); err != nil {
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
		userFound, err = service.userRepo.FindByEmail(email)
	} else {
		userFound, err = service.userRepo.FindByUserName(userName)
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
	*user = *userFound
	return nil
}

func (s *sUserService) Update(id uuid.UUID, updateMap map[string]any) (*model.User, error) {
	for key := range updateMap {
		if _, exists := model.AllowedUserFieldTags[key]; !exists {
			return nil, util.NewBadRequest(fmt.Sprintf("Invalid user info: '%s' to update", key))
		}
	}
	if password, exists := updateMap["password"]; exists {
		password := password.(string)
		password, err := util.EncodePassword(password)
		if err != nil {
			return nil, util.NewInternal("Could not encode the password")
		}
	}
	return s.userRepo.Update(id, updateMap)
}

func (s *sUserService) UpdateProfileImage(id uuid.UUID, imageFile multipart.File, imageType string) (imageURL string, err error) {
	newURL, err := s.profileImageRepo.Update(id, imageFile, imageType)
	if err != nil {
		return "", err
	}

	err = s.userRepo.UpdateProfileImageURL(id, newURL)
	if err != nil {
		return "", err
	}
	return newURL, nil
}
