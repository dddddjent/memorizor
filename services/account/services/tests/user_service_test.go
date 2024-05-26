package services_test

import (
	"memorizor/services/account/model"
	"memorizor/services/account/repository/mocks"
	"memorizor/services/account/services"
	"memorizor/services/account/util"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetByUUID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id, _ := uuid.NewV4()
		user := &model.User{
			UUID: id,
			Name: "AAA",
		}
		mockRepo := &repository.SMockUserRepository{}
		mockRepo.On("FindByUUID", id).Return(user, nil)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})

		expectUser := user
		actualUser, actualErr := service.GetByUUID(id)
		assert.Equal(t, expectUser, actualUser)
		assert.Nil(t, actualErr)
		mockRepo.AssertExpectations(t)
	})
	t.Run("NotFound", func(t *testing.T) {
		id, _ := uuid.NewV4()
		mockRepo := &repository.SMockUserRepository{}
		err := &util.Error{}
		mockRepo.On("FindByUUID", id).Return(nil, err)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})

		actualUser, actualErr := service.GetByUUID(id)
		assert.Nil(t, actualUser)
		assert.Error(t, actualErr)
		mockRepo.AssertExpectations(t)
	})
}

func TestSignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id, _ := uuid.NewV4()
		user := &model.User{
			UUID: id,
			Name: "AAA",
		}
		mockRepo := &repository.SMockUserRepository{}
		mockRepo.On("Create", user).Return(nil)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})
		err := service.SignUp(user)
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Conflict", func(t *testing.T) {
		expectErr := &util.Error{
			Type:    util.ConflictError,
			Message: "Conflict",
		}
		mockRepo := &repository.SMockUserRepository{}
		mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(expectErr)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})
		err := service.SignUp(&model.User{})
		assert.Equal(t, expectErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSignIn(t *testing.T) {
	t.Run("Successfully sign in by user name", func(t *testing.T) {
		user := &model.User{
			UserName: "AAA",
			Email:    "",
			Password: "123456",
		}
		encodedPassword, err := util.EncodePassword(user.Password)
		userFound := &model.User{
			Password: encodedPassword,
		}
		assert.NoError(t, err)
		mockRepo := &repository.SMockUserRepository{}
		mockRepo.On("FindByUserName", user.UserName).Return(userFound, nil)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})
		err = service.SignIn(user)
		assert.Nil(t, err)
		assert.Equal(t, user, userFound)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "FindByEmail")
	})
	t.Run("Successfully sign in by email", func(t *testing.T) {
		user := &model.User{
			UserName: "",
			Email:    "130604@g.com",
			Password: "123456",
		}
		encodedPassword, err := util.EncodePassword(user.Password)
		userFound := &model.User{
			Password: encodedPassword,
		}
		assert.NoError(t, err)
		mockRepo := &repository.SMockUserRepository{}
		mockRepo.On("FindByEmail", user.Email).Return(userFound, nil)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})
		err = service.SignIn(user)
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "FindByUserName")
	})
	t.Run("Incorrect email and user name combo", func(t *testing.T) {
		user := &model.User{
			UserName: "",
			Email:    "",
			Password: "123456",
		}
		service := services.NewSUserService(&services.SUserServiceConfig{})
		err := service.SignIn(user)
		assert.Error(t, err)

		assert.Equal(t, util.ErrorType("BAD_REQUEST"), err.(*util.Error).Type)
	})
	t.Run("No user found", func(t *testing.T) {
		user := &model.User{
			UserName: "",
			Email:    "130604@g.com",
			Password: "123456",
		}
		expectErr := util.NewAuthorization("No user found")
		mockRepo := &repository.SMockUserRepository{}
		mockRepo.On("FindByEmail", user.Email).Return(nil, expectErr)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})
		err := service.SignIn(user)
		assert.Equal(t, expectErr, err.(*util.Error))
		mockRepo.AssertExpectations(t)
	})
	t.Run("Incorrect password", func(t *testing.T) {
		user := &model.User{
			UserName: "",
			Email:    "130604@g.com",
			Password: "123456",
		}
		encodedPassword, err := util.EncodePassword("123455")
		userFound := &model.User{
			Password: encodedPassword,
		}
		assert.NoError(t, err)
		mockRepo := &repository.SMockUserRepository{}
		mockRepo.On("FindByEmail", user.Email).Return(userFound, nil)

		service := services.NewSUserService(&services.SUserServiceConfig{
			Repository: mockRepo,
		})
		err = service.SignIn(user)
		expectErr := util.NewAuthorization("Incorrect password")
		assert.Equal(t, expectErr, err.(*util.Error))
		mockRepo.AssertExpectations(t)
	})
}
