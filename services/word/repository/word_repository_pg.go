package repository

import (
	"fmt"
	"log"
	"memorizor/services/word/model"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type sWordRepository struct {
	db *gorm.DB
}

func NewSWordRepositoryPG(db *gorm.DB) IWordRepository {
	return &sWordRepository{db: db}
}

func (r *sWordRepository) AllWords(userID uuid.UUID, method string, offset int64, pageLength int64) []model.WordCard {
	words := []model.WordCard{}
	r.db.Where("user_id = ?", userID).
		Order(fmt.Sprintf("%s ASC", method)).
		Offset(int(offset)).Limit(int(pageLength)).
		Find(&words)
	return words
}

func (r *sWordRepository) CountAllWords(userID uuid.UUID) int64 {
	wordCnt := int64(0)
	r.db.Model(&model.WordCard{}).Where("user_id = ?", userID).Count(&wordCnt)
	log.Println("Words: ", wordCnt)
	return wordCnt
}
