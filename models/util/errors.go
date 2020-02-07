package util

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

// Define common errors
var (
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
