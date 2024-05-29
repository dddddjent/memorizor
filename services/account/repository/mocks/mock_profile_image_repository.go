package repository

import (
	"mime/multipart"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type SMockProfileImageRepository struct {
	mock.Mock
}

func (r *SMockProfileImageRepository) Update(userID uuid.UUID, imageFile multipart.File, imageType string) (imageURL string, err error) {
	args := r.Called(userID, imageFile, imageType)
	return args.String(0), args.Error(1)
}
