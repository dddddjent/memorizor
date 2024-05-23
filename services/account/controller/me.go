package controller

import (
	"memorizor/services/account/http_err"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func (ctrl *Controller) me(c *gin.Context) {
	var id uuid.UUID
	if err := id.Parse(c.Query("uuid")); err != nil {
		err := httpErr.Error{Type: httpErr.BadRequest, Message: "Can't parse uuid"}
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	user, err := ctrl.userService.GetByUUID(id)
	if err != nil {
		err, _ := err.(*httpErr.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
