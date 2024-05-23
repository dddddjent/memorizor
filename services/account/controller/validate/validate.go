package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	httpErr "memorizor/services/account/http_err"
)

type InvalidArg struct {
	Field string
	Value string
	Tag   string
	Param string
}

func ShouldBindOrBadRequest(c *gin.Context, data any) bool {
	if err := c.ShouldBind(data); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			invalidArgs := []InvalidArg{}

			for _, err := range errs {
				invalidArgs = append(invalidArgs, InvalidArg{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}
			err := httpErr.Error{
				Type:    httpErr.BadRequest,
				Message: "See invalid args for details",
			}
			c.JSON(err.HttpStatus(), gin.H{
				"error":        err,
				"invalid_args": invalidArgs,
			})
			return false
		}
        
		err := httpErr.Error{
			Type:    httpErr.Internal,
			Message: "Internal error",
		}
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return false
	}
	return true
}
