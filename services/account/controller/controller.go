package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	router *gin.Engine
}

type Config struct {
	Router *gin.Engine
}

func NewController(config *Config) *Controller {
	ctrl := &Controller{router: config.Router}
	group := ctrl.router.Group(os.Getenv("ACCOUNT_API_URL"))

	group.GET("/", ctrl.root)
	group.GET("/me", ctrl.me)
	group.POST("/signup", ctrl.signup)
	group.POST("/signin", ctrl.signin)
	group.POST("/signout", ctrl.signout)
	group.POST("/tokens", ctrl.tokens)
	group.POST("/image", ctrl.image)
	group.POST("/details", ctrl.details)

	return ctrl
}

func (ctrl *Controller) Run(addr string) {
	ctrl.router.Run(addr)
}

func (ctrl *Controller) root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (ctrl *Controller) me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "me",
	})
}

func (ctrl *Controller) signup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "signup",
	})
}

func (ctrl *Controller) signin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "signin",
	})
}

func (ctrl *Controller) signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "signout",
	})
}

func (ctrl *Controller) tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "tokens",
	})
}

func (ctrl *Controller) image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "image",
	})
}

func (ctrl *Controller) details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "details",
	})
}
