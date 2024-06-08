package middleware

import (
	"bytes"
	"io"
	"log"
	"memorizor/services/word/controller/validate"
	"memorizor/services/word/services"
	"memorizor/services/word/util"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AccessToken string `header:"Authorization"`
}

func AuthUser(s services.ITokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := authHeader{}
		if validate.ShouldBindHeaderOrBadRequest(ctx, &header) == false {
			ctx.Abort()
			return
		}

		byteBody, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(byteBody))
		rawTokenString := header.AccessToken
		tokenString := rawTokenString[7:]
		log.Println(tokenString)
		user, err := s.ValidateAccessToken(tokenString)
		if err != nil {
			err := err.(*util.Error)
			ctx.JSON(err.HttpStatus(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
