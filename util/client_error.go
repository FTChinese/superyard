package util

import "github.com/go-sql-driver/mysql"

// UnprocessableCode is an enum for UnprocessableError's Code field
type UnprocessableCode string

const (
	// CodeMissing means a resource does not exist
	CodeMissing UnprocessableCode = "missing"
	// CodeMissingField means a required field on a resource has not been set.
	CodeMissingField UnprocessableCode = "missing_field"
	// CodeInvalid means the formatting of a field is invalid
	CodeInvalid UnprocessableCode = "invalid"
	// CodeAlreadyExsits means another resource has the same value as this field.
	CodeAlreadyExsits UnprocessableCode = "already_exists"
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

// ClientError respond to 4xx http status.
type ClientError struct {
	Message string  `json:"message"`
	Reason  *Reason `json:"error,omitempty"`
}

// Reason tells why client request errored.
type Reason struct {
	message string
	Field   string            `json:"field"`
	Code    UnprocessableCode `json:"code"`
}

// NewReason creates a new instance of Reason
func NewReason() *Reason {
	return &Reason{message: "Validation failed"}
}

// NewReasonAlreadyExists creates a Reason instance with Code set to already_exists
func NewReasonAlreadyExists(field string) *Reason {
	return &Reason{
		message: "Validation failed",
		Field:   field,
		Code:    CodeAlreadyExsits,
	}
}

// SetMessage set the message to be carried away.
func (r *Reason) SetMessage(msg string) {
	r.message = msg
}

// GetMessage returns Reason's descriptive message.
func (r *Reason) GetMessage() string {
	return r.message
}
