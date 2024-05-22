package main

import (
	"log"
	"memorizor/services/account/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Account service begins.")

	r := gin.Default()
	ctrl := controller.NewController(&controller.Config{Router: r})

	ctrl.Run(":8080")
}
