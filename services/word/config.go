package main

import (
	"log"
	"memorizor/services/word/controller"
	"memorizor/services/word/services"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateDataSources() map[string]any {
    return make(map[string]any)
}

func ConfigureRouter(r *gin.Engine) {
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
	idTimeout, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_TIMEOUT"), 10, 64)
	tokenService := services.NewSTokenService(&services.STokenServiceConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		AccessTokenTimeout:  idTimeout,
	})

	requestTimeout, _ := strconv.ParseInt(os.Getenv("REQUEST_TIMEOUT"), 10, 64)
	controller.NewController(&controller.Config{
		Router:       r,
		BaseURL:      os.Getenv("WORD_API_URL"),
		Timeout:      requestTimeout,
        TokenService: tokenService,
	})
}
