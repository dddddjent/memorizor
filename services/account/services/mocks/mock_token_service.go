package services

import (
	"memorizor/services/account/model"

	"github.com/stretchr/testify/mock"
)

type SMockTokenService struct {
	mock.Mock
}

func (s *SMockTokenService) CreatePairFromUser(user *model.User, prevToken string) (*model.TokenPair, error) {
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
			panic("Can't cast arg 1 to err")
		}
		return token, err
	}
	return token, nil
}
