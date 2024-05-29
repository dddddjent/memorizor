package repository

import (
	"memorizor/services/account/model"
	"memorizor/services/account/util"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type sUserRepositoryPG struct {
	db *gorm.DB
}

func NewSUserRepositoryPG(db *gorm.DB) IUserRepository {
	return &sUserRepositoryPG{db}
}

func (r *sUserRepositoryPG) Create(user *model.User) error {
	cnt := int64(0)
	r.db.Model(&model.User{}).Where("user_name= ?", user.UserName).Count(&cnt)
	if cnt == 1 {
		return util.NewConflict("user_name", user.UserName)
	}
	r.db.Create(user)
	// put uuid into it
	r.db.Model(&model.User{}).Where("user_name= ?", user.UserName).First(user)
	return nil
}

func (r *sUserRepositoryPG) FindByUUID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	r.db.Where("uuid = ?", id).First(user)
	if *user == *new(model.User) {
		return nil, util.NewNotFound("uuid", id.String())
	}
	return user, nil
}

func (r *sUserRepositoryPG) FindByUserName(userName string) (*model.User, error) {
	user := &model.User{}
	r.db.Where("user_name = ?", userName).First(user)
	if *user == *new(model.User) {
		return nil, util.NewNotFound("user_name", userName)
	}
	return user, nil
}

func (r *sUserRepositoryPG) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	r.db.Where("email = ?", email).First(user)
	if *user == *new(model.User) {
		return nil, util.NewNotFound("email", email)
	}
	return user, nil
}

func (r *sUserRepositoryPG) Update(id uuid.UUID, update_map map[string]any) (*model.User, error) {
	r.db.Model(&model.User{}).Where("uuid = ?", id).Updates(update_map)
	return r.FindByUUID(id)
}

func (r *sUserRepositoryPG) UpdateProfileImageURL(id uuid.UUID, newURL string) error {
	user, err := r.FindByUUID(id)
	if err != nil {
		return err
	}
	r.db.Model(user).Where("uuid = ?", id).Update("profile_image_url", newURL)
	return nil
}
