package controller

import (
	"fmt"
	"log"
	"memorizor/services/account/controller/validate"
	"memorizor/services/account/model"
	"memorizor/services/account/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type signoutBody struct {
	TokenString string `json:"refresh_token"`
}

func (ctrl *sController) signout(c *gin.Context) {
	body := signoutBody{}
	if validate.ShouldBindJSONOrBadRequest(c, &body) == false {
		return
	}
	refreshToken := &model.SRefreshToken{
		ID: uuid.Nil,
	}
	if body.TokenString != "" {
		err := fmt.Errorf("")
		if refreshToken, err = ctrl.tokenService.ValidateRefreshToken(body.TokenString); err != nil {
			err := err.(*util.Error)
			c.JSON(err.HttpStatus(), gin.H{
				"error": err,
			})
			return
		}
	}

	userAny, exists := c.Get("user")
	if !exists {
		err := util.NewBadRequest("No user info found in the request")
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}
	user := userAny.(*model.User)
	log.Println(user)

	if err := ctrl.tokenService.SignOut(user, refreshToken.ID); err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sign out successfully!",
	})
}
