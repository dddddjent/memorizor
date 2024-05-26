package controller

import (
	"memorizor/services/account/controller/validate"
	"memorizor/services/account/model"
	"memorizor/services/account/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signInBody struct {
	UserName string `json:"user_name" binding:"lte=30"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,gte=6,lte=10"`
}

func (ctrl *sController) signin(c *gin.Context) {
	body := &signInBody{}
	if validate.ShouldBindJSONOrBadRequest(c, body) == false {
		return
	}

	user := &model.User{
		UserName: body.UserName,
		Email:    body.Email,
		Password: body.Password,
	}

	if err := ctrl.userService.SignIn(user); err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}

	var tokenPair *model.TokenPair
	tokenPair, err := ctrl.tokenService.CreatePairFromUser(user, "")
	if err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token_pair": tokenPair,
	})
}
