package repository

import (
	"memorizor/services/account/model"
	"memorizor/services/account/util"
	"os"

	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SUserRepositoryPG struct {
	db *gorm.DB
}

func NewSUserRepositoryPG() *SUserRepositoryPG {
	dsn :=
		"host=postgres-users user=" + os.Getenv("POSTGRES_USER") +
			" password=" + os.Getenv("POSTGRES_PASSWORD") +
			" dbname=administrator port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can't connect to Postgres")
	}
	db.AutoMigrate(&model.User{})
	return &SUserRepositoryPG{db}
}

func (r *SUserRepositoryPG) Create(user *model.User) error {
	cnt := int64(1)
	r.db.Where("user_name = ?", user.UserName).Count(&cnt)
	if cnt != 1 {
		return &util.Error{
			Type:    util.Conflict,
			Message: "Duplicate UserName",
		}
	}
	r.db.Create(user)
	return nil
}

func (r *SUserRepositoryPG) FindByUUID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	r.db.Where("uuid = ?", id.String()).First(user)
	if *user == *new(model.User) {
		return nil, &util.Error{
			Type:    util.NotFound,
			Message: "No user found",
		}
	}
	return user, nil
}
