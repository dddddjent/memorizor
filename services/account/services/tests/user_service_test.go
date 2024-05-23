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
	t.Run("Not created", func(t *testing.T) {
		expectErr := &util.Error{
			Type:    util.Conflict,
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
