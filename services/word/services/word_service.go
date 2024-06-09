package services

import (
	"fmt"
	"log"
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

func fibPreviousDays(fib int) (start time.Time, end time.Time) {
	now := time.Now()
	then := now.AddDate(0, 0, -fib)
	start = time.Date(then.Year(), then.Month(), then.Day(), 0, 0, 0, 0, then.Location())
	end = start.AddDate(0, 0, 1)
	return
}

func (s *sWordService) WordsToday(userID uuid.UUID) ([][]model.Word, error) {
	oldestDate, err := s.wordRepo.OldestCreatedTime(userID)
	if err != nil {
		return make([][]model.Word, 0), nil
	}
    log.Println("Oldest: ", oldestDate)

	i := 0
	fibNum := [2]int{0, 1}
	wordsEachDay := make([][]model.Word, 0)
	for {
		fibNum[i%2] = fibNum[i%2] + fibNum[(i+1)%2]
		start, end := fibPreviousDays(fibNum[i%2])
		log.Println("start: ", start)
		log.Println("end: ", end)
		if oldestDate.Compare(end) != -1 {
			break
		}
		words, err := s.wordRepo.WordsInRange(userID, start, end)
		if err != nil {
			return wordsEachDay, util.NewInternal("Could not get the words in this time range")
		}
		wordsEachDay = append(wordsEachDay, words)
		i++
	}
	return wordsEachDay, nil
}
