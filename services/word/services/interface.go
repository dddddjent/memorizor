package services

import "memorizor/services/word/model"

type ITokenService interface {
	ValidateAccessToken(tokenString string) (*model.User, error)
}
