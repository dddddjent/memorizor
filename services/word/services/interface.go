package services

import (
	"memorizor/services/word/model"

	"github.com/gofrs/uuid"
)

type ITokenService interface {
	ValidateAccessToken(tokenString string) (*model.User, error)
}

type IWordService interface {
	AllWords(userID uuid.UUID, method string, page int64) ([]model.WordCard, error)
}
