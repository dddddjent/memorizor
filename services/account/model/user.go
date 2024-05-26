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

func (u *User) DeepCopyTo(dst *User) {
	dst.CreatedAt = u.CreatedAt
	dst.UpdatedAt = u.UpdatedAt
	dst.UUID = u.UUID
	dst.Name = u.Name
	dst.UserName = u.UserName
	dst.Password = u.Password
	dst.Email = u.Email
	dst.ImageURL = u.ImageURL
	dst.Website = u.Website
}
