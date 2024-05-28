package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"memorizor/services/account/controller/validate"
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

type updateMeBody struct {
	Name     *string `json:"name" binding:"omitempty,gte=1,lte=40"`
	UserName *string `json:"user_name" binding:"omitempty,gte=3,lte=30"`
	Password *string `json:"password" binding:"omitempty,gte=6,lte=10"`
	Email    *string `json:"email" binding:"omitempty,email"`
	ImageURL string  `json:"image_url"`
	Website  string  `json:"website"`
	Bio      string  `json:"bio"`
}

func (ctrl *sController) updateMe(c *gin.Context) {
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

	// Only for format check
	byteBody, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteBody))
	if validate.ShouldBindJSONOrBadRequest(c, &updateMeBody{}) == false {
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteBody))
	raw, err := c.GetRawData()
	if err != nil {
		err := util.NewBadRequest("Unable to get request data")
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	updateMap := make(map[string]any)
	err = json.Unmarshal(raw, &updateMap)
	if err != nil {
		log.Println("Unable to get update map")
		err := util.NewBadRequest("Unable to get update map")
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}

	user, err = ctrl.userService.Update(id, updateMap)
	if err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
