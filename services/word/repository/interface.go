package repository

import (
	"memorizor/services/word/model"

	"github.com/gofrs/uuid"
)

type IWordRepository interface {
	AllWords(userID uuid.UUID, method string, offset int64, pageLength int64) []model.WordCard
	CountAllWords(userID uuid.UUID) int64
}
