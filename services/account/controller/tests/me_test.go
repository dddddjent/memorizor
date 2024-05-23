package controller_test

import (
	"encoding/json"
	"memorizor/services/account/controller"
	"memorizor/services/account/util"
	"memorizor/services/account/model"
	"memorizor/services/account/services/mocks"
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
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		request, _ := http.NewRequest(http.MethodGet, baseURL+"/me?uuid="+id.String(), nil)
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
		// mock the service
		userService := &services.SMockUserService{}
		userService.On("GetByUUID", "1").Return(nil, nil)

		r := gin.Default()
		controller.NewController(&controller.Config{ // r on /me
			Router:      r,
			UserService: userService,
		})

		request, _ := http.NewRequest(http.MethodGet, baseURL+"/me?uuid=1", nil)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		err := util.Error{Type: util.BadRequest, Message: "Can't parse uuid"}
		expectResponseBody, _ := json.Marshal(gin.H{
			"error": err,
		})
		expectCode := err.HttpStatus()

		assert.Equal(t, expectCode, recorder.Code)
		assert.Equal(t, expectResponseBody, recorder.Body.Bytes())
		userService.AssertNotCalled(t, "GetByUUID")
	})
	t.Run("NotFound", func(t *testing.T) {
		// mock the service
		id, _ := uuid.NewV4()
		userService := &services.SMockUserService{}
		err := &util.Error{Type: util.NotFound, Message: "Can't find the user"}
		userService.On("GetByUUID", id).Return(nil, err)

		r := gin.Default()
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
