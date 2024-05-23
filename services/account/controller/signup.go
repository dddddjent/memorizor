package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) signup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "signup",
	})
}
