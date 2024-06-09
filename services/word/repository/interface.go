package repository

import (
	"memorizor/services/word/model"
	"time"

	"github.com/gofrs/uuid"
)

type IWordRepository interface {
	AllWords(userID uuid.UUID, method string, offset int64, pageLength int64) []model.Word
	CountAllWords(userID uuid.UUID) int64
	SetWord(userID uuid.UUID, word *model.Word) error
	DeleteWord(userID uuid.UUID, wordID uuid.UUID) error
	UpdateClickedAt(userID uuid.UUID, wordID uuid.UUID, newTime time.Time) error
    // OldestWord(userID uuid.UUID) (time.Time, error)
}
