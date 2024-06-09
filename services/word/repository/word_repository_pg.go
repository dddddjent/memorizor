package repository

import (
	"fmt"
	"log"
	"memorizor/services/word/model"
	"memorizor/services/word/util"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type sWordRepository struct {
	db *gorm.DB
}

func NewSWordRepositoryPG(db *gorm.DB) IWordRepository {
	return &sWordRepository{db: db}
}

func (r *sWordRepository) AllWords(userID uuid.UUID, method string, offset int64, pageLength int64) []model.Word {
	words := []model.Word{}
	r.db.Where("user_id = ?", userID).
		Order(fmt.Sprintf("%s ASC", method)).
		Offset(int(offset)).Limit(int(pageLength)).
		Find(&words)
	return words
}

func (r *sWordRepository) CountAllWords(userID uuid.UUID) int64 {
	wordCnt := int64(0)
	r.db.Model(&model.Word{}).Where("user_id = ?", userID).Count(&wordCnt)
	log.Println("Words: ", wordCnt)
	return wordCnt
}

func (r *sWordRepository) SetWord(userID uuid.UUID, word *model.Word) error {
	if word == nil {
		return util.NewInternal("word is nil")
	}

	count := int64(0)
	r.db.Model(word).Where("user_id = ?", userID).Where("word = ?", word.Word).Count(&count)
	if count != 0 {
		r.db.Model(word).Where("word = ?", word.Word).Updates(map[string]any{
			"explanation": word.Explanation,
			"url":         word.URL,
		})
	} else {
		word.UserID = userID
		r.db.Create(word)
		r.db.Where("word = ?", word.Word).First(word)
		r.db.Model(word).Where("word = ?", word.Word).Update("clicked_at", word.CreatedAt)
	}
	return nil
}

func (r *sWordRepository) DeleteWord(userID uuid.UUID, wordID uuid.UUID) error {
	count := int64(0)
	r.db.Model(&model.Word{}).Where("id = ?", wordID).Count(&count)
	if count == 0 {
		return util.NewBadRequest("Could not find this word")
	}
	r.db.Delete(&model.Word{}, wordID)
	return nil
}

func (r *sWordRepository) UpdateClickedAt(userID uuid.UUID, wordID uuid.UUID, newTime time.Time) error {
	count := int64(0)
	r.db.Model(&model.Word{}).Where("id = ?", wordID).Count(&count)
	if count == 0 {
		return util.NewBadRequest("Could not find this word")
	}
	r.db.Model(&model.Word{}).Where("id = ?", wordID).Update("clicked_at", newTime)
	return nil
}
