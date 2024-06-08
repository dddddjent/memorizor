package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type WordCard struct {
	CreatedAt   time.Time `json:"created_at"`
	ID          uuid.UUID `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid();primarykey"`
	Word        string    `json:"word" gorm:"column:word"`
	Explanation string    `json:"explanation" gorm:"column:explanation"`
	URL         string    `json:"url" gorm:"column:url"`
	ClickedAt   time.Time `json:"clicked_at" gorm:"column:clicked_at"`
	UserID      uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid"`
	User        User      `json:"-" gorm:"references:UUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
