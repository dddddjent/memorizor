package services

import (
	"fmt"
	"memorizor/services/word/model"
	"memorizor/services/word/repository"
	"memorizor/services/word/util"
	"time"

	"github.com/gofrs/uuid"
)

type sWordService struct {
	wordRepo   repository.IWordRepository
	pageLength int64
}

type SWordServiceConfig struct {
	WordRepository repository.IWordRepository
	PageLength     int64
}

func NewSWordService(config *SWordServiceConfig) IWordService {
	return &sWordService{
		wordRepo:   config.WordRepository,
		pageLength: config.PageLength,
	}
}

func (s *sWordService) AllWords(userID uuid.UUID, method string, page int64) ([]model.Word, error) {
	if page < 1 {
		return []model.Word{}, util.NewBadRequest(fmt.Sprintf("Invalid page number %d", page))
	}

	var orderMethod string
	switch method {
	case "time":
		orderMethod = "created_at"
	case "alphabetic":
		orderMethod = "word"
	default:
		return []model.Word{}, util.NewBadRequest("Incorrect ordering method")
	}

	offset := (page - 1) * s.pageLength
	return s.wordRepo.AllWords(userID, orderMethod, offset, s.pageLength), nil
}

func (s *sWordService) CountPage(userID uuid.UUID) (pageCnt int64, err error) {
	pageCnt = (s.wordRepo.CountAllWords(userID))/int64(s.pageLength) + 1
	return pageCnt, nil
}

func (s *sWordService) SetWord(userID uuid.UUID, word *model.Word) error {
	return s.wordRepo.SetWord(userID, word)
}

func (s *sWordService) DeleteWord(userID uuid.UUID, wordID uuid.UUID) error {
	return s.wordRepo.DeleteWord(userID, wordID)
}

func (s *sWordService) ClickWord(userID uuid.UUID, wordID uuid.UUID) error {
	return s.wordRepo.UpdateClickedAt(userID, wordID, time.Now())
}
