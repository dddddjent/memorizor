package services

import (
	"memorizor/services/account/model"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type SMockUserService struct {
	mock.Mock
}

func (s *SMockUserService) GetByUUID(id uuid.UUID) (*model.User, error) {
	args := s.Called(id) // if is called by this id

	arg0 := args.Get(0) // then return these args (definded by On() method)
	var user *model.User
	if arg0 != nil {
		user = arg0.(*model.User)
	}
	err := args.Get(1)
	if err != nil {
		err, ok := err.(error)
		if !ok {
			panic("Could not cast arg 1 to err")
		}
		return user, err
	}
	return user, nil
}

func (s *SMockUserService) SignUp(user *model.User) error {
	args := s.Called(user)

	return args.Error(0)
}

func (s *SMockUserService) SignIn(user *model.User) error {
	args := s.Called(user)

	return args.Error(0)
}

func (s *SMockUserService) Update(id uuid.UUID, update_map map[string]any) (*model.User, error) {
	args := s.Called(id, update_map) // if is called by this id

	arg0 := args.Get(0) // then return these args (definded by On() method)
	var user *model.User
	if arg0 != nil {
		user = arg0.(*model.User)
	}
	err := args.Get(1)
	if err != nil {
		err, ok := err.(error)
		if !ok {
			panic("Could not cast arg 1 to err")
		}
		return user, err
	}
	return user, nil
}
