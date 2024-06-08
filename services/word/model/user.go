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
	UUID            uuid.UUID `json:"uuid" gorm:"column:uuid;type:uuid;default:gen_random_uuid();primarykey"`
	Name            string    `json:"name" gorm:"column:name"`
	UserName        string    `json:"user_name" gorm:"column:user_name"`
	Password        string    `json:"-" gorm:"column:password"`
	Email           string    `json:"email" gorm:"column:email;unique"`
	ProfileImageURL string    `json:"profile_image_url" gorm:"column:profile_image_url"`
	Website         string    `json:"website" gorm:"column:website"`
	Bio             string    `json:"bio" gorm:"column:bio"`
}
