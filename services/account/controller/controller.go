package controller

import (
	"memorizor/services/account/services"
	"net/http"

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
}

func NewController(config *Config) *sController {
	ctrl := &sController{
		router:       config.Router,
		userService:  config.UserService,
		tokenService: config.TokenService,
	}
	group := ctrl.router.Group(config.BaseURL)

	group.GET("/me", ctrl.me)
	group.POST("/signup", ctrl.signup)
	group.POST("/signin", ctrl.signin)
	group.POST("/signout", ctrl.signout)
	group.POST("/tokens", ctrl.tokens)
	group.POST("/image", ctrl.image)
	group.POST("/details", ctrl.details)

	return ctrl
}

func (ctrl *sController) Run(addr string) {
	ctrl.router.Run(addr)
}

func (ctrl *sController) signin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "signin",
	})
}

func (ctrl *sController) signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "signout",
	})
}

func (ctrl *sController) tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "tokens",
	})
}

func (ctrl *sController) image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "image",
	})
}

func (ctrl *sController) details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "details",
	})
}
