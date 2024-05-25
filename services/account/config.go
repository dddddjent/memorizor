package main

import (
	"log"
	"memorizor/services/account/controller"
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/services"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func generateDataSources() map[string]any {
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DATABASE") +
		" port=" + os.Getenv("POSTGRES_PORT")
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can't connect to Postgres")
	}
	postgresDB.AutoMigrate(&model.User{})

	return map[string]any{
		"postgres": postgresDB,
	}
}

func ConfigureRouter(r *gin.Engine) {
	dataSources := generateDataSources()

	userService := services.NewSUserService(&services.SUserServiceConfig{
		Repository: repository.NewSUserRepositoryPG(dataSources["postgres"].(*gorm.DB)),
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
	idTimeout, err := strconv.ParseInt(os.Getenv("ID_TOKEN_TIMEOUT"), 10, 64)
	refreshTimeout, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_TIMEOUT"), 10, 64)
	tokenService := services.NewSTokenService(&services.STokenServiceConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		RefreshSecret:       refreshSecret,
		IdTokenTimeout:      idTimeout,
		RefreshTokenTimeout: refreshTimeout,
	})

	controller.NewController(&controller.Config{
		Router:       r,
		UserService:  userService,
		TokenService: tokenService,
		BaseURL:      os.Getenv("ACCOUNT_API_URL"),
	})
}
