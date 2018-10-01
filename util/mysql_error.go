package util

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// DuplicateError carries message for MySQL insertion duplicate error
// See https://dev.mysql.com/doc/refman/8.0/en/error-messages-server.html `Error: 1062`
// Message: Duplicate entry '%s' for key %d
type DuplicateError struct {
	Field   string
	Number  uint16
	Message string
}

func (de DuplicateError) Error() string {
	return fmt.Sprintf("Error %d: %s", de.Number, de.Message)
}

// SQLInsertError wraps MySQLError to DuplicateError is it is a duplicate error, or return the original error otherwise.
func SQLInsertError(err error, field string) error {
	if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return DuplicateError{
			Field:   field,
			Number:  e.Number,
			Message: e.Message,
		}
	}

	return err
}
