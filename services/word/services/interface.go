package services

import (
	"memorizor/services/word/model"

	"github.com/gofrs/uuid"
)

type ITokenService interface {
	ValidateAccessToken(tokenString string) (*model.User, error)
}

type IWordService interface {
	AllWords(userID uuid.UUID, method string, page int64) ([]model.Word, error)
	CountPage(userID uuid.UUID) (pageCnt int64, err error)
	SetWord(userID uuid.UUID, word *model.Word) error
	DeleteWord(userID uuid.UUID, wordID uuid.UUID) error
	ClickWord(userID uuid.UUID, wordID uuid.UUID) error
}
