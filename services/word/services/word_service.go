package services

import (
	"fmt"
	"memorizor/services/word/model"
	"memorizor/services/word/repository"
	"memorizor/services/word/util"

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

func (s *sWordService) AllWords(userID uuid.UUID, method string, page int64) ([]model.WordCard, error) {
	if page < 1 {
		return []model.WordCard{}, util.NewBadRequest(fmt.Sprintf("Invalid page number %d", page))
	}

	var orderMethod string
	switch method {
	case "time":
		orderMethod = "created_at"
	case "alphabetic":
		orderMethod = "word"
	default:
		return []model.WordCard{}, util.NewBadRequest("Incorrect ordering method")
	}
	// pageCnt := (s.wordRepo.CountAllWords(userID))/int64(s.pageLength) + 1
	offset := (page - 1) * s.pageLength
	return s.wordRepo.AllWords(userID, orderMethod, offset, s.pageLength), nil
}
