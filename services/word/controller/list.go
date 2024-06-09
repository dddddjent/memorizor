package controller

import (
	"log"
	"memorizor/services/word/model"
	"memorizor/services/word/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ctrl *sController) list(c *gin.Context) {
	page, err := strconv.ParseInt(c.Param("page"), 10, 64)
	if err != nil {
		util.ResponseDefaultError(c, util.NewBadRequest(err.Error()))
		return
	}
	log.Println("Page: ", page)
	method := c.Query("method")
	log.Println("Method: ", method)

	userAny, exists := c.Get("user")
	if !exists {
		util.ResponseDefaultError(c, util.NewBadRequest("No user info found in the request"))
		return
	}
	user := userAny.(*model.User)
	id := user.UUID
	log.Println("ID: ", id)

	wordList, err := ctrl.wordService.AllWords(id, method, page)
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list": wordList,
	})
}
