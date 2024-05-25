package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Account service begins.")

	r := gin.Default()
	ConfigureRouter(r)

	r.Run(":8080")
}
