package main

import (
	"log"
	"memorizor/services/word/controller"
	"memorizor/services/word/model"
	"memorizor/services/word/repository"
	"memorizor/services/word/services"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func generateDataSources() map[string]any {
	dsn := "host=" + os.Getenv("POSTGRES_USERS_HOST") +
		" user=" + os.Getenv("POSTGRES_USERS_USER") +
		" password=" + os.Getenv("POSTGRES_USERS_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_USERS_DATABASE") +
		" port=" + os.Getenv("POSTGRES_USERS_PORT")
	log.Println("Connecting to Postgres.")
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to Postgres")
	}
	postgresDB.AutoMigrate(&model.WordCard{})

	return map[string]any{
		"postgres": postgresDB,
	}
}

func ConfigureRouter(r *gin.Engine) {
	dataSources := generateDataSources()

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
		PrivateKey:         privateKey,
		PublicKey:          publicKey,
		AccessTokenTimeout: idTimeout,
	})

	pageLength, err := strconv.ParseInt(os.Getenv("PAGE_LENGTH"), 0, 64)
	if err != nil {
		log.Fatal(err.Error())
	}
	wordRepository := repository.NewSWordRepositoryPG(dataSources["postgres"].(*gorm.DB))
	wordService := services.NewSWordService(&services.SWordServiceConfig{
		WordRepository: wordRepository,
		PageLength:     pageLength,
	})

	requestTimeout, _ := strconv.ParseInt(os.Getenv("REQUEST_TIMEOUT"), 10, 64)
	controller.NewController(&controller.Config{
		Router:       r,
		BaseURL:      os.Getenv("WORD_API_URL"),
		Timeout:      requestTimeout,
		TokenService: tokenService,
		WordService:  wordService,
	})
}
