package main

import (
	"context"
	"log"
	"memorizor/services/account/controller"
	"memorizor/services/account/model"
	"memorizor/services/account/repository"
	"memorizor/services/account/services"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
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
	postgresDB.AutoMigrate(&model.User{})

	log.Println("Connecting to Redis.")
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("Could not connect to Redis")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles(
			[]string{"/keys" + os.Getenv("AWS_CREDENTIAL_PATH")},
		),
	)
	if err != nil {
		panic("Could not connect to aws")
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = os.Getenv("AWS_REGION")
	})

	return map[string]any{
		"postgres":            postgresDB,
		"redis":               rdb,
		"aws":                 client,
		"aws_bucket_name":     os.Getenv("AWS_BUCKET_NAME"),
		"aws_bucket_url_root": os.Getenv("AWS_BUCKET_URL_ROOT"),
	}
}

func ConfigureRouter(r *gin.Engine) {
	dataSources := generateDataSources()

	tokenRepository := repository.NewSTokenRepositoryRedis(dataSources["redis"].(*redis.Client))

	userService := services.NewSUserService(&services.SUserServiceConfig{
		UserRepository: repository.NewSUserRepositoryPG(dataSources["postgres"].(*gorm.DB)),
		ProfileImageRepo: repository.NewSProfileImageRepositoryAWS(
			dataSources["aws"].(*s3.Client),
			dataSources["aws_bucket_name"].(string),
			dataSources["aws_bucket_url_root"].(string),
		),
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
	idTimeout, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_TIMEOUT"), 10, 64)
	refreshTimeout, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_TIMEOUT"), 10, 64)
	tokenService := services.NewSTokenService(&services.STokenServiceConfig{
		TokenRepository:     tokenRepository,
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		RefreshSecret:       refreshSecret,
		AccessTokenTimeout:  idTimeout,
		RefreshTokenTimeout: refreshTimeout,
	})

	requestTimeout, _ := strconv.ParseInt(os.Getenv("REQUEST_TIMEOUT"), 10, 64)
	controller.NewController(&controller.Config{
		Router:       r,
		UserService:  userService,
		TokenService: tokenService,
		BaseURL:      os.Getenv("ACCOUNT_API_URL"),
		Timeout:      requestTimeout,
	})
}
