package util_test

import (
	"memorizor/services/account/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	t.Run("Correct Password", func(t *testing.T) {
		password := "123"
		encoded, err := util.EncodePassword(password)
		assert.NoError(t, err)

		result, err := util.ComparePassword(encoded, password)
		assert.NoError(t, err)

		assert.True(t, result)
	})
	t.Run("Incorrect Password", func(t *testing.T) {
		password := "123"
		passwordFalse := "123sd"
		encoded, err := util.EncodePassword(password)
		assert.NoError(t, err)

		result, err := util.ComparePassword(encoded, passwordFalse)
		assert.NoError(t, err)

		assert.False(t, result)
	})
}
