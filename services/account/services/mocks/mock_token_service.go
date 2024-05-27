package services

import (
	"memorizor/services/account/model"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type SMockTokenService struct {
	mock.Mock
}

func (s *SMockTokenService) CreatePairFromUser(user *model.User, prevToken uuid.UUID) (*model.TokenPair, error) {
	args := s.Called(user, prevToken)

	arg0 := args.Get(0)
	var token *model.TokenPair
	if arg0 != nil {
		token = arg0.(*model.TokenPair)
	}
	err := args.Get(1)
	if err != nil {
		err, ok := err.(error)
		if !ok {
			panic("Could not cast arg 1 to err")
		}
		return token, err
	}
	return token, nil
}

func (s *SMockTokenService) ValidateAccessToken(tokenString string) (*model.User, error) {
	args := s.Called(tokenString)

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

func (s *SMockTokenService) ValidateRefreshToken(tokenString string) (*model.SRefreshToken, error) {
	args := s.Called(tokenString)

	arg0 := args.Get(0)
	var token *model.SRefreshToken
	if arg0 != nil {
		token = arg0.(*model.SRefreshToken)
	}
	err := args.Get(1)
	if err != nil {
		err, ok := err.(error)
		if !ok {
			panic("Could not cast arg 1 to err")
		}
		return token, err
	}
	return token, nil
}
