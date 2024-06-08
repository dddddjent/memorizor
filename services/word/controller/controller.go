package controller

import (
	"memorizor/services/word/controller/middleware"
	"memorizor/services/word/model"
	"memorizor/services/word/services"
	"memorizor/services/word/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type sController struct {
	router       *gin.Engine
	tokenService services.ITokenService
}

type Config struct {
	Router       *gin.Engine
	BaseURL      string
	Timeout      int64
	TokenService services.ITokenService
}

func NewController(config *Config) *sController {
	ctrl := &sController{
		router: config.Router,
        tokenService: config.TokenService,
	}

	rootGroup := ctrl.router.Group(config.BaseURL)
	timeoutDuration := time.Duration(config.Timeout) * time.Second
	if gin.Mode() != gin.TestMode {
		rootGroup.Use(middleware.Timeout(timeoutDuration, util.NewServiceUnavailable()))
		rootGroup.Use(middleware.AuthUser(ctrl.tokenService))
	}
	{
		rootGroup.GET("/page", func(c *gin.Context) {
			userAny, exists := c.Get("user")
			if !exists {
				err := util.NewBadRequest("No user info found in the request")
				c.JSON(err.HttpStatus(), gin.H{
					"error": err,
				})
				return
			}
			user := userAny.(*model.User)
			c.JSON(http.StatusOK, gin.H{
				"user": user,
			})
		})
	}

	return ctrl
}
