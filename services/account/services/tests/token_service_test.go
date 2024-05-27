package services_test

import (
	"fmt"
	"log"
	"memorizor/services/account/model"
	mockRepo "memorizor/services/account/repository/mocks"
	"memorizor/services/account/services"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePairFromUser(t *testing.T) {
	privateKeyBytes, err := os.ReadFile("../../keys/rsa_private_test.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	publicKeyBytes, err := os.ReadFile("../../keys/rsa_public_test.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	secret := "A secret"
	idTimeOut := 900
	refreshTimeOut := 259200

	id, _ := uuid.NewV4()
	user := &model.User{
		UUID:     id,
		UserName: "me",
	}
	id, _ = uuid.NewV4()
	userForError := &model.User{
		UUID:     id,
		UserName: "me",
	}

	mockRepository := &mockRepo.SMockTokenRepository{}
	mockRepository.On("SetRefreshToken", user.UUID, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockRepository.On("SetRefreshToken", userForError.UUID, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("time.Duration")).Return(fmt.Errorf("An error"))
	mockRepository.On("DeleteRefreshToken", user.UUID, mock.AnythingOfType("uuid.UUID")).Return(nil)
	tokenService := services.NewSTokenService(&services.STokenServiceConfig{
		TokenRepository:     mockRepository,
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		RefreshSecret:       secret,
		AccessTokenTimeout:  int64(idTimeOut),
		RefreshTokenTimeout: int64(refreshTimeOut),
	})

	type accessTokenClaims struct {
		jwt.RegisteredClaims
		User *model.User `json:"user"`
	}
	type refreshTokenClaims struct {
		jwt.RegisteredClaims
		UUID    string `json:"uuid"`
		TokenId string `json:"token_id"`
	}
	t.Run("Generate a pair of keys", func(t *testing.T) {
		tokenPair, err := tokenService.CreatePairFromUser(user, uuid.Nil)
		assert.NoError(t, err)

		// AccessToken
		accessClaims := accessTokenClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.AccessToken, &accessClaims, func(t *jwt.Token) (any, error) {
			return publicKey, nil
		})
		assert.NoError(t, err)

		actualUser := accessClaims.User

		// same user
		assert.Equal(t, user, actualUser)
		assert.Empty(t, actualUser.Password)

		// time issue
		actualExpire := time.Unix(accessClaims.ExpiresAt.Unix(), 0)
		expectedExpire := time.Now().Add(time.Duration(idTimeOut) * time.Second)
		assert.WithinDuration(t, actualExpire, expectedExpire, 5*time.Second) // this and create pair should be within 5s

		// RefreshToken
		refreshClaims := refreshTokenClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken.TokenString, &refreshClaims, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		assert.NoError(t, err)

		actualUUIDString := refreshClaims.UUID
		assert.Equal(t, actualUUIDString, user.UUID.String())

		actualExpire = time.Unix(refreshClaims.ExpiresAt.Unix(), 0)
		expectedExpire = time.Now().Add(time.Duration(refreshTimeOut) * time.Second)
		assert.WithinDuration(t, actualExpire, expectedExpire, 5*time.Second) // this and create pair should be within 5s
	})

	t.Run("Unable to set refresh token", func(t *testing.T) {
		_, err := tokenService.CreatePairFromUser(userForError, uuid.Nil)
		assert.Error(t, err)
		mockRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Unable to delete refresh token", func(t *testing.T) {
		mockRepository.On("DeleteRefreshToken", userForError.UUID, mock.AnythingOfType("uuid.UUID")).Return(fmt.Errorf("Not good"))
		id, err := uuid.NewV4()
		_, err = tokenService.CreatePairFromUser(userForError, id)
		assert.Error(t, err)
		mockRepository.AssertNotCalled(t, "SetRefreshToken")
	})

	t.Run("Empty previous token id", func(t *testing.T) {
		_, err := tokenService.CreatePairFromUser(user, uuid.Nil)
		assert.NoError(t, err)
		mockRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})
}

func TestValidateTokens(t *testing.T) {
	privateKeyBytes, err := os.ReadFile("../../keys/rsa_private_test.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	publicKeyBytes, err := os.ReadFile("../../keys/rsa_public_test.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	secret := "A secret"
	idTimeOut := 1
	refreshTimeOut := 1

	id, _ := uuid.NewV4()
	user := &model.User{
		UUID:     id,
		UserName: "me",
	}

	mockRepository := &mockRepo.SMockTokenRepository{}
	mockRepository.On("SetRefreshToken", user.UUID, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockRepository.On("DeleteRefreshToken", user.UUID, mock.AnythingOfType("uuid.UUID")).Return(nil)
	tokenService := services.NewSTokenService(&services.STokenServiceConfig{
		TokenRepository:     mockRepository,
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		RefreshSecret:       secret,
		AccessTokenTimeout:  int64(idTimeOut),
		RefreshTokenTimeout: int64(refreshTimeOut),
	})
	t.Run("Success", func(t *testing.T) {
		tokenPair, err := tokenService.CreatePairFromUser(user, uuid.Nil)
		assert.NoError(t, err)

		actualUser, err := tokenService.ValidateAccessToken(tokenPair.AccessToken)
		assert.NoError(t, err)
		refreshToken, err := tokenService.ValidateRefreshToken(tokenPair.RefreshToken.TokenString)
		assert.NoError(t, err)

		assert.Equal(t, actualUser, user)
		assert.Equal(t, user.UUID, refreshToken.UUID)
	})
	t.Run("Broken token", func(t *testing.T) {
		tokenPair, err := tokenService.CreatePairFromUser(user, uuid.Nil)
		assert.NoError(t, err)

		_, err = tokenService.ValidateAccessToken(tokenPair.AccessToken[1:])
		assert.Error(t, err)

		_, err = tokenService.ValidateRefreshToken(tokenPair.RefreshToken.TokenString[1:])
		assert.Error(t, err)
	})
	t.Run("Expired", func(t *testing.T) {
		tokenPair, err := tokenService.CreatePairFromUser(user, uuid.Nil)
		assert.NoError(t, err)

		time.Sleep(time.Duration(1) * time.Second)

		_, err = tokenService.ValidateAccessToken(tokenPair.AccessToken)
		assert.Error(t, err)

		_, err = tokenService.ValidateRefreshToken(tokenPair.RefreshToken.TokenString[1:])
		assert.Error(t, err)
	})
}
