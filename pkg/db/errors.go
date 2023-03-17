package db

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// IsAlreadyExists tests if an error means the field already exists
func IsAlreadyExists(err error) bool {
	if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return true
	}

	return false
}

// ConvertGormError tests if an error is gorm's
// ErrRecordNotFound and change to sql.ErrNoRows if it is
// to keep API output unchanged.
func ConvertGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return sql.ErrNoRows
	}

	return err
}
