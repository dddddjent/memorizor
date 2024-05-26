package validate

import (
	"fmt"
	"memorizor/services/account/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InvalidArg struct {
	Field string
	Value string
	Tag   string
	Param string
}

func ShouldBindJSONOrBadRequest(c *gin.Context, data any) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts application/json", c.FullPath())
		err := util.NewUnsupportedMediaType(msg)
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return false
	}
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
			err := util.Error{
				Type:    util.BadRequestError,
				Message: "See invalid args for details",
			}
			c.JSON(err.HttpStatus(), gin.H{
				"error":        err,
				"invalid_args": invalidArgs,
			})
			return false
		}

		err := util.Error{
			Type:    util.InternalError,
			Message: "Internal error",
		}
		c.JSON(err.HttpStatus(), gin.H{
			"error": err,
		})
		return false
	}
	return true
}
