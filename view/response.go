package view

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Response collects all data needed for an HTTP response
type Response struct {
	StatusCode int
	Header     http.Header
	Body       interface{}
	IsError    bool
}

// SetBody sets reponse body to any value
func (r Response) SetBody(b interface{}) Response {
	r.Body = b
	return r
}

// NoCache set headers to prevent the response from being cached.
func (r Response) NoCache() Response {
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Cache-Control", "no-store")
	r.Header.Add("Cache-Control", "must-revalidate")
	r.Header.Add("Pragma", "no-cache")
	return r
}

// NewResponse creates a new instance of Response with default values
func NewResponse() Response {
	r := Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		IsError:    false,
	}

	r.Header.Set("Content-Type", "application/json; charset=utf-8")

	return r
}

// NewNoContent creates an HTTP 204 No Content response
func NewNoContent() Response {
	r := NewResponse()
	r.StatusCode = http.StatusNoContent

	return r
}

// NewNotFound creates response 404 Not Found
func NewNotFound() Response {
	r := NewResponse()

	r.StatusCode = http.StatusNotFound
	r.Body = ClientError{Message: "Not Found."}

	return r
}

// NewUnauthorized create a new instance of Response for 401 Unauthorized response
func NewUnauthorized(msg string) Response {
	if msg == "" {
		msg = "Requires authorization."
	}

	r := NewResponse()
	r.StatusCode = http.StatusUnauthorized
	r.Body = ClientError{Message: msg}

	return r
}

// NewForbidden creates response for StatusForbidden
func NewForbidden(msg string) Response {
	r := NewResponse()

	r.StatusCode = http.StatusForbidden
	r.Body = ClientError{Message: msg}

	return r
}

// NewBadRequest creates a new Response for 400 Bad Request with the specified msg
func NewBadRequest(msg string) Response {

	if msg == "" {
		msg = "Problems parsing JSON"
	}
	r := NewResponse()

	r.StatusCode = http.StatusBadRequest
	r.Body = ClientError{Message: msg}

	return r
}

// NewUnprocessable creates response 422 Unprocessable Entity
func NewUnprocessable(msg string, err error) Response {
	if msg == "" {
		msg = "Validation Failed"
	}

	ce := ClientError{Message: msg}
	if err != nil {
		ce.ErrorDetail = err
	}

	r := NewResponse()
	r.StatusCode = http.StatusUnprocessableEntity
	r.Body = ce

	return r
}

// NewInternalError creates response for internal server error
func NewInternalError(msg string) Response {

	r := NewResponse()

	r.StatusCode = http.StatusInternalServerError
	r.Body = ClientError{Message: msg}

	return r
}

// NewDBFailure handles MySQl errors
func NewDBFailure(err error) Response {

	switch err := err.(type) {
	case util.DuplicateError:
		// field: "email", code: "already_exists"
		ue := UnprocessableError{
			Field: err.Field,
			Code:  CodeAlreadyExsits.String(),
		}
		return NewUnprocessable(err.Message, ue)
	}

	switch err {
	case sql.ErrNoRows:
		return NewNotFound()

	default:
		return NewInternalError(err.Error())
	}
}
