package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type GormModel struct {
	ID        uint         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" gorm:"index"`
}

type User struct {
	GormModel
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	UserName string    `json:"user_name"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	ImageURL string    `json:"image_url"`
	Website  string    `json:"website"`
}
