package controller

import (
	"memorizor/services/account/controller/middleware"
	"memorizor/services/account/services"
	"memorizor/services/account/util"
	"time"

	"github.com/gin-gonic/gin"
)

type sController struct {
	router       *gin.Engine
	userService  services.IUserService
	tokenService services.ITokenService
}

type Config struct {
	Router       *gin.Engine
	UserService  services.IUserService
	TokenService services.ITokenService
	BaseURL      string
	Timeout      int64
}

func NewController(config *Config) *sController {
	ctrl := &sController{
		router:       config.Router,
		userService:  config.UserService,
		tokenService: config.TokenService,
	}

	rootGroup := ctrl.router.Group(config.BaseURL)
	timeoutDuration := time.Duration(config.Timeout) * time.Second
	if gin.Mode() != gin.TestMode {
		rootGroup.Use(middleware.Timeout(timeoutDuration, util.NewServiceUnavailable()))
	}

	group := rootGroup.Group(".")
	{
		group.POST("/signup", ctrl.signup)
		group.POST("/signin", ctrl.signin)
		group.POST("/token", ctrl.token)
	}

	authGroup := rootGroup.Group(".")
	if gin.Mode() != gin.TestMode {
		authGroup.Use(middleware.AuthUser(ctrl.tokenService))
	}
	{
		authGroup.GET("/me", ctrl.me)
		authGroup.POST("/me", ctrl.updateMe)
		authGroup.POST("/signout", ctrl.signout)
		authGroup.POST("/profile_image", ctrl.profile_image)
	}

	return ctrl
}
