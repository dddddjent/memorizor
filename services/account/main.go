package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Account service begins.")

	r := gin.Default()

	r.GET("/api/account", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "O",
		})
	})

	r.Run(":8080")
}
