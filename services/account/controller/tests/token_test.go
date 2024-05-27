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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	id, _ := uuid.NewV4()
	user1 := &model.User{
		UUID:     id,
		Password: "123456",
	}
	id, _ = uuid.NewV4()
	user2 := &model.User{
		UUID:     id,
		Password: "123456",
	}
	id, _ = uuid.NewV4()
	user3 := &model.User{
		UUID:     id,
		Password: "123456",
	}
	secret := "A secret"
	timeout := time.Duration(1) * time.Second
	rToken1, _ := util.GenerateRefreshToken(user1.UUID, secret, timeout)
	tokenPair1 := model.TokenPair{
		RefreshToken: *rToken1,
	}
	rToken2, _ := util.GenerateRefreshToken(user2.UUID, secret, timeout)
	rToken3, _ := util.GenerateRefreshToken(user3.UUID, secret, timeout)

	mockUserService := &services.SMockUserService{}
	mockTokenService := &services.SMockTokenService{}
	baseURL := os.Getenv("ACCOUNT_API_URL")
	r := gin.Default()
	controller.NewController(&controller.Config{
		Router:       r,
		UserService:  mockUserService,
		TokenService: mockTokenService,
		Timeout:      5,
		BaseURL:      baseURL,
	})

	expectErr := util.NewBadRequest("Not found")
	mockTokenService.On("ValidateRefreshToken", rToken1.TokenString).Return(rToken1, nil)
	mockTokenService.On("ValidateRefreshToken", rToken2.TokenString).Return(rToken2, nil)
	mockTokenService.On("ValidateRefreshToken", rToken3.TokenString).Return(rToken3, nil)
	mockUserService.On("GetByUUID", rToken1.UUID).Return(user1, nil)
	mockUserService.On("GetByUUID", rToken2.UUID).Return(nil, expectErr)
	mockUserService.On("GetByUUID", rToken3.UUID).Return(user3, nil)
	mockTokenService.On("CreatePairFromUser", user1, rToken1.ID).Return(&tokenPair1, nil)
	mockTokenService.On("CreatePairFromUser", user3, rToken3.ID).Return(nil, expectErr)
	t.Run("Success", func(t *testing.T) {
		postBody, _ := json.Marshal(map[string]any{
			"refresh_token": rToken1.TokenString,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/token",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusOK
		assert.Equal(t, expectCode, recorder.Code)

		expectBody, err := json.Marshal(map[string]*model.TokenPair{"token_pair": &tokenPair1})
		assert.NoError(t, err)
		assert.Equal(t, expectBody, recorder.Body.Bytes())
	})

	t.Run("Validate failed", func(t *testing.T) {
		expectErr := util.NewAuthorization("Broken token")
		mockTokenService.On("ValidateRefreshToken", rToken1.TokenString[1:]).Return(
			nil,
			expectErr,
		)
		postBody, _ := json.Marshal(map[string]any{
			"refresh_token": rToken1.TokenString[1:],
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/token",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := expectErr.HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)

		mockUserService.AssertNotCalled(t, "GetByUUID")
	})

	t.Run("No user found", func(t *testing.T) {
		postBody, _ := json.Marshal(map[string]any{
			"refresh_token": rToken2.TokenString,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/token",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := expectErr.HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)

		mockUserService.AssertNotCalled(t, "GetByUUID")
	})

	t.Run("Unable to create a new pair", func(t *testing.T) {
		postBody, _ := json.Marshal(map[string]any{
			"refresh_token": rToken3.TokenString,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/token",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := expectErr.HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)
	})
}
