package controller

import (
	"log"
	"memorizor/services/word/model"
	"memorizor/services/word/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *sController) page(c *gin.Context) {
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
	log.Println("ID: ", id)

	pageCnt, err := ctrl.wordService.CountPage(id)
	if err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pages": pageCnt,
	})
}
