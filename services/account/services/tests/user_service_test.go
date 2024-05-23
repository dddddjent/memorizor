package services_test

import (
	httpErr "memorizor/services/account/http_err"
	"memorizor/services/account/model"
	"memorizor/services/account/repository/mocks"
	"memorizor/services/account/services"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
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
		err := &httpErr.Error{}
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
