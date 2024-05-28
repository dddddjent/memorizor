package model

import (
	"reflect"
	"time"

	"github.com/gofrs/uuid"
)

type GormModel struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	GormModel
	UUID     uuid.UUID `json:"uuid" gorm:"column:uuid;type:uuid;default:gen_random_uuid();primarykey"`
	Name     string    `json:"name" gorm:"column:name"`
	UserName string    `json:"user_name" gorm:"column:user_name"`
	Password string    `json:"-" gorm:"column:password"`
	Email    string    `json:"email" gorm:"column:email;unique"`
	ImageURL string    `json:"image_url" gorm:"column:image_url"`
	Website  string    `json:"website" gorm:"column:website"`
	Bio      string    `json:"bio" gorm:"column:bio"`
}

var AllowedUserFieldTags map[string]struct{} = func() map[string]struct{} {
	invalidTags := map[string]struct{}{
		"created_at": {},
		"updated_at": {},
		"uuid":       {},
		"-":          {},
	}
	userType := reflect.TypeOf(User{})
	tags := make(map[string]struct{})
	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		tag := field.Tag.Get("json")
		if _, exists := invalidTags[tag]; !exists {
			tags[tag] = struct{}{}
		}
	}
	tags["password"] = struct{}{}
	return tags
}()
