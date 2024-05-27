package controller

import (
	"log"
	"memorizor/services/account/controller/validate"
	"memorizor/services/account/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tokenBody struct {
	TokenString string `json:"refresh_token" binding:"required"`
}

// only for refresh token
func (ctrl *sController) token(c *gin.Context) {
	body := tokenBody{}
	if validate.ShouldBindJSONOrBadRequest(c, &body) == false {
		return
	}

	refreshToken, err := ctrl.tokenService.ValidateRefreshToken(body.TokenString)
	if err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}

	user, err := ctrl.userService.GetByUUID(refreshToken.UUID)
	log.Println(refreshToken.ID)
	if err != nil {
		err := err.(*util.Error)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return
	}

	tokenPair, err := ctrl.tokenService.CreatePairFromUser(user, refreshToken.ID)
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
