package service

import (
	"errors"
	"fmt"

	"github.com/sfarosu/go-tooling-portal/internal/helper"
)

// GeneratePassword produces a random password with the requested length and character classes.
// Returns an error when inputs are invalid (length <= 0, length too large, or no character class selected).
func GeneratePassword(length int, uppercase, lowercase, numbers, symbols bool) (string, error) {
	const maxLength = 1024

	if length <= 0 {
		return "", errors.New("length must be greater than zero")
	}
	if length > maxLength {
		return "", fmt.Errorf("length must be <= %d", maxLength)
	}
	if !uppercase && !lowercase && !numbers && !symbols {
		return "", errors.New("at least one character class must be selected")
	}

	// helper.RandomString returns a generated string based on the options
	pwd, err := helper.RandomString(length, uppercase, lowercase, numbers, symbols)
	if err != nil {
		return "", err
	}
	return pwd, nil
}
