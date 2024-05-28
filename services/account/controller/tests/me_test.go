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
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var baseURL string = os.Getenv("ACCOUNT_API_URL")
	t.Run("Success", func(t *testing.T) {
		// mock the service
		id, _ := uuid.NewV4()
		user := &model.User{
			UUID: id,
			Name: "AAA",
		}
		userService := &services.SMockUserService{}
		userService.On("GetByUUID", id).Return(user, nil)

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", user)
		})
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		request, _ := http.NewRequest(http.MethodGet, baseURL+"/me", nil)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectResponseBody, _ := json.Marshal(gin.H{
			"user": user,
		})
		expectCode := http.StatusOK

		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, expectResponseBody, recorder.Body.Bytes())
		userService.AssertExpectations(t)
	})
	t.Run("BadRequest", func(t *testing.T) {
		id, _ := uuid.NewV4()
		user := &model.User{
			UUID: id,
			Name: "AAA",
		}
		userService := &services.SMockUserService{}
		userService.On("GetByUUID", user.UUID).Return(nil, nil)

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			// ctx.Set("user", user)
		})
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		request, _ := http.NewRequest(http.MethodGet, baseURL+"/me", nil)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		err := util.NewBadRequest("No user info found in the request")
		expectResponseBody, _ := json.Marshal(gin.H{
			"error": err,
		})
		expectCode := err.HttpStatus()

		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, expectResponseBody, recorder.Body.Bytes())
		userService.AssertNotCalled(t, "GetByUUID")
	})
	t.Run("NotFound", func(t *testing.T) {
		id, _ := uuid.NewV4()
		user := &model.User{
			UUID: id,
			Name: "AAA",
		}
		userService := &services.SMockUserService{}
		err := &util.Error{Type: util.NotFoundError, Message: "Could not find the user"}
		userService.On("GetByUUID", id).Return(nil, err)

		r := gin.Default()
		r.Use(func(ctx *gin.Context) {
			ctx.Set("user", user)
		})
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		request, _ := http.NewRequest(http.MethodGet, baseURL+"/me?uuid="+id.String(), nil)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectResponseBody, _ := json.Marshal(gin.H{
			"error": err,
		})
		expectCode := err.HttpStatus()

		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, expectResponseBody, recorder.Body.Bytes())
		userService.AssertExpectations(t)
	})
}

func TestUpdateMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var baseURL string = os.Getenv("ACCOUNT_API_URL")

	id, _ := uuid.NewV4()
	user := &model.User{
		UUID: id,
		Name: "AAA",
	}
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("user", user)
	})
	updateMap := map[string]any{
		"user_name": "123345",
		"email":     "1@a.com",
		"Password":  "123456",
	}
	userService := &services.SMockUserService{}
	userService.On("Update", id, updateMap).Return(user, nil)
	controller.NewController(&controller.Config{ // r on /me
		Router:      r,
		UserService: userService,
	})
	t.Run("Correct body", func(t *testing.T) {
		postBody, _ := json.Marshal(updateMap)
		request, err := http.NewRequest(http.MethodPost, baseURL+"/me", bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := http.StatusOK
		assert.Equal(t, expectCode, recorder.Code)
	})
	t.Run("Invalid email", func(t *testing.T) {
		updateMap := map[string]any{
			"email": "12",
		}
		postBody, _ := json.Marshal(updateMap)
		request, err := http.NewRequest(http.MethodPost, baseURL+"/me", bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := util.NewBadRequest("").HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)
	})
	t.Run("Invalid user name", func(t *testing.T) {
		updateMap := map[string]any{
			"user_name": "1ldfjalsdfjlasdjfalsdjflasdjflaksdjflajsdlf",
		}
		postBody, _ := json.Marshal(updateMap)
		request, err := http.NewRequest(http.MethodPost, baseURL+"/me", bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := util.NewBadRequest("").HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)
	})
	t.Run("Invalid name", func(t *testing.T) {
		updateMap := map[string]any{
			"name": "",
		}
		postBody, _ := json.Marshal(updateMap)
		request, err := http.NewRequest(http.MethodPost, baseURL+"/me", bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := util.NewBadRequest("").HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)
	})
	t.Run("Invalid password", func(t *testing.T) {
		updateMap := map[string]any{
			"password": "123",
		}
		postBody, _ := json.Marshal(updateMap)
		request, err := http.NewRequest(http.MethodPost, baseURL+"/me", bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		expectCode := util.NewBadRequest("").HttpStatus()
		assert.Equal(t, expectCode, recorder.Code)
	})
}
