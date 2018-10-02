package util

import (
	"database/sql"
	"errors"
	"net/http"
)

// Flags returned by some function to tell caller what kind of error response should be used
var (
	ErrBadRequest = errors.New("response: bad request")
)

// Response collects all data needed for an HTTP response
type Response struct {
	StatusCode int
	Header     http.Header
	Body       interface{}
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
	}

	r.Header.Set("Content-Type", "application/json; charset=utf-8")

	return r
}

// NewNoContent creates an HTTP 204 No Content response
func NewNoContent() Response {
	r := NewResponse().NoCache()
	r.StatusCode = http.StatusNoContent

	return r
}

// NewNotFound creates response 404 Not Found
func NewNotFound() Response {
	r := NewResponse().NoCache()

	r.StatusCode = http.StatusNotFound
	r.Body = ClientError{Message: "Not Found"}

	return r
}

// NewUnauthorized create a new instance of Response for 401 Unauthorized response
func NewUnauthorized(msg string) Response {
	if msg == "" {
		msg = "Requires authorization."
	}

	r := NewResponse().NoCache()
	r.StatusCode = http.StatusUnauthorized
	r.Body = ClientError{Message: msg}

	return r
}

// NewForbidden creates response for 403
func NewForbidden(msg string) Response {
	r := NewResponse().NoCache()

	r.StatusCode = http.StatusForbidden
	r.Body = ClientError{Message: msg}

	return r
}

// NewBadRequest creates a new Response for 400 Bad Request with the specified msg
func NewBadRequest(msg string) Response {

	if msg == "" {
		msg = "Problems parsing JSON"
	}
	r := NewResponse().NoCache()

	r.StatusCode = http.StatusBadRequest
	r.Body = ClientError{Message: msg}

	return r
}

// NewUnprocessable creates response 422 Unprocessable Entity
func NewUnprocessable(vr ValidationResult) Response {

	if vr.Message == "" {
		vr.Message = "Validation Failed"
	}

	r := NewResponse().NoCache()
	r.StatusCode = http.StatusUnprocessableEntity
	r.Body = vr

	return r
}

// NewInternalError creates response for internal server error
func NewInternalError(msg string) Response {

	r := NewResponse().NoCache()

	r.StatusCode = http.StatusInternalServerError
	r.Body = ClientError{Message: msg}

	return r
}

// NewDBFailure handles various errors returned from the model layter
// It could DuplicateError when inserting into a uniquely constrained column;
// ErrNoRows if it cannot retrieve any rows of the specified criteria;
// Otherwise internal server error.
func NewDBFailure(err error, field string) Response {

	// if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
	// 	ue := ValidationResult{
	// 		Field: field,
	// 		Code:  CodeAlreadyExsits,
	// 	}

	// 	return NewUnprocessable(e.Message, ue)
	// }

	switch err {
	case sql.ErrNoRows:
		return NewNotFound()

	default:
		return NewInternalError(err.Error())
	}
}
