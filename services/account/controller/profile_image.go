package controller

import (
	"memorizor/services/account/model"
	"memorizor/services/account/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *sController) profile_image(c *gin.Context) {
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

	imageFileHeader, err := c.FormFile("image")
	if err != nil {
		err := util.NewBadRequest(err.Error())
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	imageTypeRaw := imageFileHeader.Header.Get("Content-Type")
	imageType, err := util.ExtractImageType(imageTypeRaw)
	if err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	imageFile, err := imageFileHeader.Open()
	if err != nil {
		err := util.NewBadRequest("Could not extract the image")
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	newURL, err := ctrl.userService.UpdateProfileImage(id, imageFile, imageType)
	c.JSON(http.StatusOK, gin.H{
		"profile_image_url": newURL,
	})
}
