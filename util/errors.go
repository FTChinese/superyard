package util

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Define common errors
var (
	ErrForbidden     = errors.New("acess denied") // Tells controller what kind of  HTTP response to use
	ErrWrongPassword = errors.New("wrong password")
	ErrAlreadyExists = errors.New("already exists")
)

// IsAlreadyExists tests if an error means the field already exists
func IsAlreadyExists(err error) bool {
	if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return true
	}

	if err == ErrAlreadyExists {
		return true
	}

	return false
}
