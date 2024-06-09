package controller

import (
	"memorizor/services/word/controller/middleware"
	"memorizor/services/word/services"
	"memorizor/services/word/util"
	"time"

	"github.com/gin-gonic/gin"
)

type sController struct {
	router       *gin.Engine
	tokenService services.ITokenService
	wordService  services.IWordService
}

type Config struct {
	Router       *gin.Engine
	BaseURL      string
	Timeout      int64
	TokenService services.ITokenService
	WordService  services.IWordService
}

func NewController(config *Config) *sController {
	ctrl := &sController{
		router:       config.Router,
		tokenService: config.TokenService,
		wordService:  config.WordService,
	}

	rootGroup := ctrl.router.Group(config.BaseURL)
	timeoutDuration := time.Duration(config.Timeout) * time.Second
	if gin.Mode() != gin.TestMode {
		rootGroup.Use(middleware.Timeout(timeoutDuration, util.NewServiceUnavailable()))
		rootGroup.Use(middleware.AuthUser(ctrl.tokenService))
	}
	{
		rootGroup.GET("/list/:page", ctrl.list)
		rootGroup.GET("/page", ctrl.page)
		rootGroup.POST("/word", ctrl.word)
	}

	return ctrl
}
