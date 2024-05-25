package main

import (
	"log"
	"memorizor/services/account/controller"
	"memorizor/services/account/repository"
	"memorizor/services/account/services"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ConfigureRouter(r *gin.Engine) {
	userService := services.NewSUserService(&services.SUserServiceConfig{
		Repository: repository.NewSUserRepositoryPG(),
	})

	privateKeyBytes, err := os.ReadFile("/keys/" + os.Getenv("RSA_PRIVATE_KEY_FILE"))
	if err != nil {
		log.Fatal(err.Error())
	}
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	publicKeyBytes, err := os.ReadFile("/keys/" + os.Getenv("RSA_PUBLIC_KEY_FILE"))
	if err != nil {
		log.Fatal(err.Error())
	}
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	refreshSecret := os.Getenv("HS_REFRESH_SECRET")
	tokenService := services.NewSTokenService(&services.STokenServiceConfig{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		RefreshSecret: refreshSecret,
	})

	controller.NewController(&controller.Config{
		Router:       r,
		UserService:  userService,
		TokenService: tokenService,
	})
}
