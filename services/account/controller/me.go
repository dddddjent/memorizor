package controller

import (
	"memorizor/services/account/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func (ctrl *sController) me(c *gin.Context) {
	var id uuid.UUID
	if err := id.Parse(c.Query("uuid")); err != nil {
		err := util.Error{Type: util.BadRequest, Message: "Could not parse uuid"}
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
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
