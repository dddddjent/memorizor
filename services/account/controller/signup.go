package controller

import (
	"memorizor/services/account/controller/validate"
	"memorizor/services/account/model"
	"memorizor/services/account/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signUpBody struct {
	UserName string `json:"user_name" binding:"required,gte=3,lte=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=10"`
}

func (ctrl *Controller) signup(c *gin.Context) {
	body := &signUpBody{}
	if validate.ShouldBindOrBadRequest(c, body) == false {
		return
	}

	user := &model.User{
		UserName: body.UserName,
		Email:    body.Email,
		Password: body.Password,
	}

	if err := ctrl.userService.SignUp(user); err != nil {
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

	c.JSON(http.StatusCreated, gin.H{
		"token_pair": tokenPair,
	})
}
