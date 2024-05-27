package repository

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type SMockTokenRepository struct {
	mock.Mock
}

func (s *SMockTokenRepository) SetRefreshToken(userID, tokenID uuid.UUID, expiresIn time.Duration) error {
	ret := s.Called(userID, tokenID, expiresIn)

	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(error)
}
func (s *SMockTokenRepository) DeleteRefreshToken(userID, previousTokenID uuid.UUID) error {
	ret := s.Called(userID, previousTokenID)

	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(error)
}
