package repository

import (
	"memorizor/services/account/model"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type SMockUserRepository struct {
	mock.Mock
}

func (r *SMockUserRepository) FindByUUID(id uuid.UUID) (*model.User, error) {
	args := r.Called(id) // if is called by this id

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

func (r *SMockUserRepository) Create(user *model.User) error {
	args := r.Called(user) // if is called by this id
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (r *SMockUserRepository) FindByUserName(userName string) (*model.User, error) {
	args := r.Called(userName)

	arg0 := args.Get(0)
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

func (r *SMockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := r.Called(email)

	arg0 := args.Get(0)
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

func (r *SMockUserRepository) Update(id uuid.UUID, update_map map[string]any) (*model.User, error) {
	args := r.Called(id, update_map)

	arg0 := args.Get(0)
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
