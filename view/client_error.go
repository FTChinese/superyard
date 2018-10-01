package view

import "fmt"

type UnprocessableCode int

const (
	CodeMissing       UnprocessableCode = 0
	CodeMissingField  UnprocessableCode = 1
	CodeInvalid       UnprocessableCode = 2
	CodeAlreadyExsits UnprocessableCode = 3
)

func (c UnprocessableCode) String() string {
	codes := [...]string{
		"missing",
		"missing_field",
		"invalid",
		"already_exists",
	}

	return codes[c]
}

// ClientError respond to 4xx http status.
type ClientError struct {
	Message     string `json:"message"`
	ErrorDetail error  `json:"error,omitempty"`
}

func (e ClientError) Error() string {
	return e.Message
}

// UnprocessableError respond to 422 status code
type UnprocessableError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
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
