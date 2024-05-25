package repository

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type SMockTokenRepository struct {
	mock.Mock
}

func (s *SMockTokenRepository) SetRefreshToken(userID, tokenID string, expiresIn time.Duration) error {
	ret := s.Called(userID, tokenID, expiresIn)

	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(error)
}
func (s *SMockTokenRepository) DeleteRefreshToken(userID, previousTokenID string) error {
	ret := s.Called(userID, previousTokenID)

	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(error)
}
