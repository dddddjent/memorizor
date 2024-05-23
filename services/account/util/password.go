package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func EncodePassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	hashPassword := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))
	return hashPassword, nil
}

func ComparePassword(store, given string) (bool, error) {
	pwsalt := strings.Split(store, ".")
	salt, err := hex.DecodeString(pwsalt[1])
	if err != nil {
		return false, fmt.Errorf("Unable to decode salt")
	}

	shash, err := scrypt.Key([]byte(given), salt, 32768, 8, 1, 32)
	return hex.EncodeToString(shash) == pwsalt[0], nil
}
