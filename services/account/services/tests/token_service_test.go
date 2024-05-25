package services_test

import (
	"encoding/json"
	"log"
	"memorizor/services/account/model"
	"memorizor/services/account/services"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreatePairFromUser(t *testing.T) {
	privateKeyBytes, err := os.ReadFile("../../../../keys/rsa_private_test.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	publicKeyBytes, err := os.ReadFile("../../../../keys/rsa_public_test.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	secret := "A secret"
	idTimeOut := 900
	refreshTimeOut := 259200

	tokenService := services.NewSTokenService(&services.STokenServiceConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		RefreshSecret:       secret,
		IdTokenTimeout:      int64(idTimeOut),
		RefreshTokenTimeout: int64(refreshTimeOut),
	})

	id, _ := uuid.NewV4()
	user := &model.User{
		UUID:     id,
		UserName: "me",
	}

	t.Run("Generate a pair of keys", func(t *testing.T) {
		tokenPair, err := tokenService.CreatePairFromUser(user, "")
		assert.NoError(t, err)

		// IDToken
		parsedClaims := make(jwt.MapClaims)
		_, err = jwt.ParseWithClaims(tokenPair.IDToken, &parsedClaims, func(t *jwt.Token) (any, error) {
			return publicKey, nil
		})
		assert.NoError(t, err)

		parsedUserBytes, err := json.Marshal(parsedClaims["user"])
		assert.NoError(t, err)
		actualUser := &model.User{}
		err = json.Unmarshal(parsedUserBytes, actualUser)
		assert.NoError(t, err)

		// same user
		assert.Equal(t, user, actualUser)
		assert.Empty(t, actualUser.Password)

		// time issue
		actualExpire := time.Unix(int64(parsedClaims["exp"].(float64)), 0)
		expectedExpire := time.Now().Add(time.Duration(idTimeOut) * time.Second)
		assert.WithinDuration(t, actualExpire, expectedExpire, 5*time.Second) // this and create pair should be within 5s

		// RefreshToken
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken, &parsedClaims, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		assert.NoError(t, err)

		actualUUIDString := parsedClaims["uuid"].(string)
		_, ok := parsedClaims["token_id"].(string)
		assert.True(t, ok)
		assert.Equal(t, actualUUIDString, user.UUID.String())

		actualExpire = time.Unix(int64(parsedClaims["exp"].(float64)), 0)
		expectedExpire = time.Now().Add(time.Duration(refreshTimeOut) * time.Second)
		assert.WithinDuration(t, actualExpire, expectedExpire, 5*time.Second) // this and create pair should be within 5s
	})
}
