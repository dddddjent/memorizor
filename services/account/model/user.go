package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type GormModel struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	GormModel
	UUID     uuid.UUID `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();primarykey"`
	Name     string    `json:"name"`
	UserName string    `json:"user_name"`
	Password string    `json:"-"`
	Email    string    `json:"email" gorm:"unique"`
	ImageURL string    `json:"image_url"`
	Website  string    `json:"website"`
}
