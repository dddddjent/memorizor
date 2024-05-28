package controller_test

import (
	"bytes"
	"encoding/json"
	"memorizor/services/account/controller"
	"memorizor/services/account/model"
	services "memorizor/services/account/services/mocks"
	"memorizor/services/account/util"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignOut(t *testing.T) {
	gin.SetMode(gin.TestMode)
	baseURL := os.Getenv("ACCOUNT_API_URL")

	id, _ := uuid.NewV4()
	user := &model.User{
		UUID:     id,
		UserName: "AAAAA",
		Email:    "333@g.com",
		Password: "123456",
	}
	t.Run("Success: single signout", func(t *testing.T) {
		tokenString := "111"
		refreshToken := &model.SRefreshToken{
			TokenString: tokenString,
		}
		mockTokenService := &services.SMockTokenService{}
		mockTokenService.On("ValidateRefreshToken", tokenString).Return(refreshToken, nil)
		mockTokenService.On("SignOut", user, mock.AnythingOfType("uuid.UUID")).Return(nil)

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", user)
		})
		controller.NewController(&controller.Config{ // r on /me
			Router:       r,
			TokenService: mockTokenService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"refresh_token": tokenString,
		})
		request, err := http.NewRequest(http.MethodPost, baseURL+"/signout", bytes.NewBuffer(postBody))
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusOK
		assert.Equal(t, expectCode, recorder.Code)
	})

	t.Run("Success: all signout", func(t *testing.T) {
		mockTokenService := &services.SMockTokenService{}
		mockTokenService.On("SignOut", user, uuid.Nil).Return(nil)

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", user)
		})
		controller.NewController(&controller.Config{ // r on /me
			Router:       r,
			TokenService: mockTokenService,
		})

		postBody, _ := json.Marshal(map[string]string{})
		request, err := http.NewRequest(http.MethodPost, baseURL+"/signout", bytes.NewBuffer(postBody))
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusOK
		assert.Equal(t, expectCode, recorder.Code)
		mockTokenService.AssertNotCalled(t, "ValidateRefreshToken")
	})

	t.Run("Failed to validate", func(t *testing.T) {
		tokenString := "111"
		expectErr := util.NewAuthorization("No")
		mockTokenService := &services.SMockTokenService{}
		mockTokenService.On("ValidateRefreshToken", tokenString).Return(nil, expectErr)

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", user)
		})
		controller.NewController(&controller.Config{ // r on /me
			Router:       r,
			TokenService: mockTokenService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"refresh_token": tokenString,
		})
		request, err := http.NewRequest(http.MethodPost, baseURL+"/signout", bytes.NewBuffer(postBody))
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := expectErr.HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)
		mockTokenService.AssertNotCalled(t, "SignOut")
	})

	t.Run("Failed to signout", func(t *testing.T) {
		tokenString := "111"
		expectErr := util.NewInternal("No")
		refreshToken := &model.SRefreshToken{
			TokenString: tokenString,
		}
		mockTokenService := &services.SMockTokenService{}
		mockTokenService.On("ValidateRefreshToken", tokenString).Return(refreshToken, nil)
		mockTokenService.On("SignOut", user, refreshToken.ID).Return(expectErr)

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", user)
		})
		controller.NewController(&controller.Config{ // r on /me
			Router:       r,
			TokenService: mockTokenService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"refresh_token": tokenString,
		})
		request, err := http.NewRequest(http.MethodPost, baseURL+"/signout", bytes.NewBuffer(postBody))
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := expectErr.HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)
		mockTokenService.AssertExpectations(t)
	})
}
