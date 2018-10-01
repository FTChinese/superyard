package view

import "fmt"

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

// ClientError respond to 4xx http status.
type ClientError struct {
	Message string `json:"message"`
	Reason  error  `json:"error,omitempty"`
}

func (e ClientError) Error() string {
	return e.Message
}

// UnprocessableError respond to 422 status code
type UnprocessableError struct {
	Resource string            `json:"resource"`
	Field    string            `json:"field"`
	Code     UnprocessableCode `json:"code"`
}

// NewUnprocessableError creates a new instance of UnprocessableError
// func NewUnprocessableError() UnprocessableError {
// 	return UnprocessableError{}
// }

func (e UnprocessableError) Error() string {
	return fmt.Sprintf("Error occured at resource - %s, field - %s, code %s", e.Resource, e.Field, e.Code)
}

// func (e UnprocessableError) SetCodeMissing() UnprocessableError {
// 	e.Code = "missing"
// 	return e
// }

// func (e UnprocessableError) SetCodeMissingField() UnprocessableError {
// 	e.Code = "missing_field"
// 	return e
// }

// func (e UnprocessableError) SetCodeInvalid() UnprocessableError {
// 	e.Code = "invalid"
// 	return e
// }

// func (e UnprocessableError) SetCodeAlreadyExists() UnprocessableError {
// 	e.Code = "already_exists"
// 	return e
// }
