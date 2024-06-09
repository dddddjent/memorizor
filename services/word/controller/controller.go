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
		rootGroup.GET("/page", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"word_card": model.WordCard{},
			})
		})
		rootGroup.GET("/list/:page", ctrl.list)
	}

	return ctrl
}
