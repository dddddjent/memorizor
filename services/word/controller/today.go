package controller

import (
	"log"
	"memorizor/services/word/model"
	"memorizor/services/word/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *sController) today(c *gin.Context) {
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

	wordsToday, err := ctrl.wordService.WordsToday(id)
	if err != nil {
		util.ResponseDefaultError(c, err)
	}
	c.JSON(http.StatusOK, gin.H{
		"today": wordsToday,
	})
}
