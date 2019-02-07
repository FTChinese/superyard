package util

import (
	"errors"
)

// Define common errors
var (
	ErrForbidden     = errors.New("access denied") // Tells controller what kind of  HTTP response to use
	ErrWrongPassword = errors.New("wrong password")
	ErrAlreadyExists = errors.New("already exists")
)
