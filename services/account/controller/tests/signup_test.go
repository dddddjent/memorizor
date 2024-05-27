package controller_test

import (
	"bytes"
	"encoding/json"
	"memorizor/services/account/controller"
	"memorizor/services/account/model"
	"memorizor/services/account/services/mocks"
	"memorizor/services/account/util"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var baseURL string = os.Getenv("ACCOUNT_API_URL")
	t.Run("Success", func(t *testing.T) {
		user := &model.User{
			UserName: "AAAAA",
			Email:    "333@g.com",
			Password: "123456",
		}
		tokenPair := &model.TokenPair{
			AccessToken:  "123",
			RefreshToken: model.SRefreshToken{},
		}
		userService := &services.SMockUserService{}
		userService.On("SignUp", user).Return(nil)
		tokenService := &services.SMockTokenService{}
		tokenService.On("CreatePairFromUser", user, "").Return(tokenPair, nil)

		r := gin.Default()
		controller.NewController(&controller.Config{ // r on /me
			Router:       r,
			UserService:  userService,
			TokenService: tokenService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"user_name": user.UserName,
			"email":     user.Email,
			"Password":  user.Password,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusCreated
		expectBody, err := json.Marshal(map[string]*model.TokenPair{"token_pair": tokenPair})
		assert.NoError(t, err)

		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, expectBody, recorder.Body.Bytes())
		userService.AssertExpectations(t)
		tokenService.AssertExpectations(t)
	})

	t.Run("User name too short or too long", func(t *testing.T) {
		user1 := &model.User{
			UserName: "A",
			Email:    "333@g.com",
			Password: "123456",
		}
		user2 := &model.User{
			UserName: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			Email:    "333@g.com",
			Password: "123456",
		}
		userService := &services.SMockUserService{}
		userService.On("SignUp", user1).Return(nil)
		userService.On("SignUp", user2).Return(nil)

		r := gin.Default()
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"user_name": user1.UserName,
			"email":     user1.Email,
			"Password":  user1.Password,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusBadRequest
		actualResp := make(map[string][]map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &actualResp)
		actualFieldErr := actualResp["invalid_args"][0]["Field"]

		assert.Equal(t, 1, len(actualResp["invalid_args"]))
		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, "UserName", actualFieldErr)
		userService.AssertNotCalled(t, "SignUp")

		postBody, _ = json.Marshal(map[string]string{
			"user_name": user2.UserName,
			"email":     user2.Email,
			"Password":  user2.Password,
		})
		request, err = http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode = http.StatusBadRequest
		err = json.Unmarshal(recorder.Body.Bytes(), &actualResp)
		actualFieldErr = actualResp["invalid_args"][0]["Field"]

		assert.Equal(t, 1, len(actualResp["invalid_args"]))
		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, "UserName", actualFieldErr)
		userService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Bad email", func(t *testing.T) {
		user1 := &model.User{
			UserName: "AAAAA",
			Email:    "333@gcom",
			Password: "123456",
		}
		user2 := &model.User{
			UserName: "AAAAA",
			Email:    "333g.com",
			Password: "123456",
		}
		userService := &services.SMockUserService{}
		userService.On("SignUp", user1).Return(nil)
		userService.On("SignUp", user2).Return(nil)

		r := gin.Default()
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"user_name": user1.UserName,
			"email":     user1.Email,
			"Password":  user1.Password,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusBadRequest
		actualResp := make(map[string][]map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &actualResp)
		actualFieldErr := actualResp["invalid_args"][0]["Field"]

		assert.Equal(t, 1, len(actualResp["invalid_args"]))
		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, "Email", actualFieldErr)
		userService.AssertNotCalled(t, "SignUp")

		postBody, _ = json.Marshal(map[string]string{
			"user_name": user2.UserName,
			"email":     user2.Email,
			"Password":  user2.Password,
		})
		request, err = http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode = http.StatusBadRequest
		err = json.Unmarshal(recorder.Body.Bytes(), &actualResp)
		actualFieldErr = actualResp["invalid_args"][0]["Field"]

		assert.Equal(t, 1, len(actualResp["invalid_args"]))
		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, "Email", actualFieldErr)
		userService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Password too long or too short", func(t *testing.T) {
		user1 := &model.User{
			UserName: "AAAAA",
			Email:    "333@g.com",
			Password: "12456",
		}
		user2 := &model.User{
			UserName: "AAAAA",
			Email:    "333@g.com",
			Password: "12345612390359",
		}
		userService := &services.SMockUserService{}
		userService.On("SignUp", user1).Return(nil)
		userService.On("SignUp", user2).Return(nil)

		r := gin.Default()
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"user_name": user1.UserName,
			"email":     user1.Email,
			"Password":  user1.Password,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusBadRequest
		actualResp := make(map[string][]map[string]string)
		err = json.Unmarshal(recorder.Body.Bytes(), &actualResp)
		actualFieldErr := actualResp["invalid_args"][0]["Field"]

		assert.Equal(t, 1, len(actualResp["invalid_args"]))
		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, "Password", actualFieldErr)
		userService.AssertNotCalled(t, "SignUp")

		postBody, _ = json.Marshal(map[string]string{
			"user_name": user2.UserName,
			"email":     user2.Email,
			"Password":  user2.Password,
		})
		request, err = http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder = httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode = http.StatusBadRequest
		err = json.Unmarshal(recorder.Body.Bytes(), &actualResp)
		actualFieldErr = actualResp["invalid_args"][0]["Field"]

		assert.Equal(t, 1, len(actualResp["invalid_args"]))
		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, "Password", actualFieldErr)
		userService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Sign up internal error", func(t *testing.T) {
		user := &model.User{
			UserName: "AAAAA",
			Email:    "333@g.com",
			Password: "123456",
		}
		userService := &services.SMockUserService{}
		expectErr := util.NewInternal("No")
		userService.On("SignUp", mock.AnythingOfType("*model.User")).Return(expectErr)

		r := gin.Default()
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"user_name": user.UserName,
			"email":     user.Email,
			"Password":  user.Password,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusInternalServerError

		assert.Equal(t, expectCode, recorder.Code)
		userService.AssertExpectations(t)
	})

	t.Run("Token internal error", func(t *testing.T) {
		user := &model.User{
			UserName: "AAAAA",
			Email:    "333@g.com",
			Password: "123456",
		}
		expectErr := util.NewInternal("No")
		userService := &services.SMockUserService{}
		userService.On("SignUp", mock.AnythingOfType("*model.User")).Return(nil)
		tokenService := &services.SMockTokenService{}
		tokenService.On("CreatePairFromUser", mock.AnythingOfType("*model.User"), "").Return(nil, expectErr)

		r := gin.Default()
		controller.NewController(&controller.Config{ // r on /me
			Router:       r,
			UserService:  userService,
			TokenService: tokenService,
		})

		postBody, _ := json.Marshal(map[string]string{
			"user_name": user.UserName,
			"email":     user.Email,
			"Password":  user.Password,
		})
		request, err := http.NewRequest(
			http.MethodPost,
			baseURL+"/signup",
			bytes.NewBuffer(postBody),
		)
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusInternalServerError

		assert.Equal(t, expectCode, recorder.Code)
		userService.AssertExpectations(t)
		tokenService.AssertExpectations(t)
	})
}
