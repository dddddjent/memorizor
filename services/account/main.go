package main

import (
	"log"
	"memorizor/services/account/controller"
	"memorizor/services/account/repository"
	"memorizor/services/account/services"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Account service begins.")

	r := gin.Default()
	// ctrl := controller.NewController(&controller.Config{
	// 	Router: r,
	// })
	ctrl := controller.NewController(&controller.Config{
		Router: r,
		UserService: services.NewSUserService(&services.SUserServiceConfig{
			Repository: repository.NewSUserRepositoryPG(),
		}),
	})

	ctrl.Run(":8080")
}
