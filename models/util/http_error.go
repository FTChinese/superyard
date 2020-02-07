package util

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/validator"
	"net/http"
)

// RestfulError is the response body returned to
// an error request.
type RestfulError struct {
	StatusCode int                      `json:"-"`
	Message    string                   `json:"message"`
	Param      string                   `json:"param,omitempty"` // Param and Type should co-exist, indicating the which filed of input data failed to pass validation and why.
	Type       validator.InputFieldCode `json:"type,omitempty"`
}

// Sets the errored field name and reason.
func (re *RestfulError) Set(p string, t validator.InputFieldCode) *RestfulError {
	re.Param = p
	re.Type = t

	return re
}

func (re *RestfulError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", re.StatusCode, re.Message)
}

func NewRestfulError(code int, message string) *RestfulError {

	return &RestfulError{
		StatusCode: code,
		Message:    message,
	}
}

func NewNotFound(msg string) *RestfulError {
	return NewRestfulError(http.StatusNotFound, msg)
}

func NewBadRequest(msg string) *RestfulError {
	return NewRestfulError(http.StatusBadRequest, msg)
}

func NewForbidden(msg string) *RestfulError {
	return NewRestfulError(http.StatusForbidden, msg)
}

func NewUnprocessable(ie *validator.InputError) *RestfulError {
	return &RestfulError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    ie.Message,
		Param:      ie.Field,
		Type:       ie.Code,
	}
}

// NewAlreadyExists creates a RestfulError for unique constraint conflict.
// This is a special case of 422 Unprocessable response.
func NewAlreadyExists(field string) *RestfulError {
	return &RestfulError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "Conflict with existing value",
		Param:      field,
		Type:       validator.CodeAlreadyExists,
	}
}

func NewDBFailure(err error) *RestfulError {
	switch err {
	case sql.ErrNoRows:
		return NewRestfulError(http.StatusNotFound, err.Error())

	default:
		return NewRestfulError(http.StatusInternalServerError, err.Error())
	}
}

func RestfulErrorHandler(err error, c echo.Context) {
	re, ok := err.(*RestfulError)
	if !ok {
		re = &RestfulError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	if re.Message == "" {
		re.Message = http.StatusText(re.StatusCode)
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(re.StatusCode)
		} else {
			err = c.JSON(re.StatusCode, re)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
