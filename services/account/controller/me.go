package controller

import (
	"memorizor/services/account/model"
	"memorizor/services/account/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *sController) me(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		err := util.NewBadRequest("No user info found in the request")
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	user := userAny.(*model.User)
	id := user.UUID
	user, err := ctrl.userService.GetByUUID(id)
	if err != nil {
		err, _ := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
