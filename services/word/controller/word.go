package controller

import (
	"fmt"
	"log"
	"memorizor/services/word/controller/validate"
	"memorizor/services/word/model"
	"memorizor/services/word/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type wordRequest struct {
	Method     string         `json:"method" binding:"required"`
	Parameters map[string]any `json:"parameters" binding:"required"`
}

func (ctrl *sController) word(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		util.ResponseDefaultError(c, util.NewBadRequest("No user info found in the request"))
		return
	}
	user := userAny.(*model.User)
	id := user.UUID
	log.Println("ID: ", id)

	body := &wordRequest{}
	if validate.ShouldBindJSONOrBadRequest(c, body) == false {
		c.Abort()
		return
	}

	switch body.Method {
	case "add", "update":
		{
			ctrl.setWord(c, id, body.Parameters)
		}
	case "delete":
		{
			ctrl.deleteWord(c, id, body.Parameters)
		}
	case "click":
		{
			ctrl.clickWord(c, id, body.Parameters)
		}
	default:
		{
			util.ResponseDefaultError(c, util.NewBadRequest("Unknown method"))
			return
		}
	}
}

func parse(param map[string]any, field string) (string, error) {
	value, exists := param[field]
	if !exists {
		err := util.NewBadRequest(field + " needs to be specified")
		return "", err
	}
	value_string, ok := value.(string)
	if !ok {
		err := util.NewBadRequest(fmt.Sprintf("Could not parse %s parameter", field))
		return "", err
	}
	return value_string, nil
}

type setParameters struct {
	Word        string `json:"word"`
	Explanation string `json:"explanation"`
	URL         string `json:"url"`
}

func formatWord(word string) (string, error) {
	if word == "" {
		return "", util.NewBadRequest("No empty word")
	}

	validChar := func(ch rune) bool {
		if ch <= 'z' && ch >= 'a' || ch >= 'A' && ch <= 'Z' {
			return true
		}
		return false
	}
	toUpper := func(ch rune) rune {
		if ch <= 'z' && ch >= 'a' {
			ch = ch - 'a' + 'A'
		}
		return ch
	}
	toLower := func(ch rune) rune {
		if ch <= 'Z' && ch >= 'A' {
			ch = ch - 'A' + 'a'
		}
		return ch
	}
	out := []rune(word)
	for i, ch := range word {
		if !validChar(ch) {
			return "", util.NewBadRequest("Invalid characters")
		}
		if i == 0 {
			out[i] = toUpper(ch)
		} else {
			out[i] = toLower(ch)
		}
	}
	return string(out), nil
}

func (ctrl *sController) setWord(c *gin.Context, userID uuid.UUID, param map[string]any) {
	err := (error)(nil)
	word := &model.Word{}
	word.Word, err = parse(param, "word")
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}
	word.Word, err = formatWord(word.Word)
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}
	word.Explanation, err = parse(param, "explanation")
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}
	word.URL, err = parse(param, "url")
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}

	if err := ctrl.wordService.SetWord(userID, word); err != nil {
		util.ResponseDefaultError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (ctrl *sController) deleteWord(c *gin.Context, userID uuid.UUID, param map[string]any) {
	wordIDString, err := parse(param, "id")
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}
	wordID := uuid.FromStringOrNil(wordIDString)
	if wordID == uuid.Nil {
		util.ResponseDefaultError(c, util.NewBadRequest("Could not parse uuid"))
		return
	}

	err = ctrl.wordService.DeleteWord(userID, wordID)
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (ctrl *sController) clickWord(c *gin.Context, userID uuid.UUID, param map[string]any) {
	wordIDString, err := parse(param, "id")
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}
	wordID := uuid.FromStringOrNil(wordIDString)
	if wordID == uuid.Nil {
		util.ResponseDefaultError(c, util.NewBadRequest("Could not parse uuid"))
		return
	}

	err = ctrl.wordService.ClickWord(userID, wordID)
	if err != nil {
		util.ResponseDefaultError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
