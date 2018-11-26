package util

import (
	"fmt"
	"unicode/utf8"

	validate "github.com/asaskevich/govalidator"
)

const (
	msgTooLong  = "The length of %s should not exceed %d chars"
	msgTooShort = "The length of %s should not less than %d chars"
	msgLenRange = "The length of %s must be within %d to %d chars"
)

// isLength tests if a string's length is within a range.
func isLength(str string, min, max int) bool {
	if min > max {
		min, max = max, min
	}
	strLength := utf8.RuneCountInString(str)
	return strLength >= min && strLength <= max
}

// minLength tests if a string's length is longer than min
func minLength(str string, min int) bool {
	strLength := utf8.RuneCountInString(str)
	return strLength >= min
}

// maxLength tests if a string's length is under max
// Return true if the length of str is under or equal to max; false otherwise
func maxLength(str string, max int) bool {
	strLength := utf8.RuneCountInString(str)
	return strLength <= max
}

// RequireNotEmpty makes sure the value is not an empty string
func RequireNotEmpty(value, field string) *Reason {
	if value == "" {
		r := NewReason()
		r.Code = CodeMissingField
		r.Field = field

		return r
	}

	return nil
}

// RequireLenRange makes sure the value's length is within the specified range
func RequireLenRange(value string, min int, max int, field string) *Reason {
	if !isLength(value, min, max) {
		r := NewReason()
		r.SetMessage(fmt.Sprintf(msgLenRange, field, min, max))
		r.Field = field
		r.Code = CodeInvalid

		return r
	}

	return nil
}

// OptionalMaxLen makes sure a string's length does not exceed the max limit.
// Empty string is valid.
func OptionalMaxLen(value string, max int, field string) *Reason {
	if !maxLength(value, max) {
		r := NewReason()
		r.SetMessage(fmt.Sprintf(msgTooLong, field, max))
		r.Field = field
		r.Code = CodeInvalid

		return r
	}

	return nil
}

// RequireNotEmptyWithMax validates a string is not empty and must not exceed max chars.
func RequireNotEmptyWithMax(value string, max int, field string) *Reason {
	if r := RequireNotEmpty(value, field); r != nil {
		return r
	}

	return OptionalMaxLen(value, max, field)
}

// RequireNotEmptyWithinLen validates a string is not empty, its length is within the specified range.
func RequireNotEmptyWithinLen(value string, min, max int, field string) *Reason {
	if r := RequireNotEmpty(value, field); r != nil {
		return r
	}

	return RequireLenRange(value, min, max, field)
}

// RequireEmail make sure the email is not empty space and is indeed an email address.
func RequireEmail(email string) *Reason {
	if r := RequireNotEmpty(email, "email"); r != nil {
		return r
	}

	if !validate.IsEmail(email) {
		r := NewReason()
		r.Code = CodeInvalid
		r.Field = "email"

		return r
	}

	return OptionalMaxLen(email, 255, "email")
}

// RequirePassword ensures the password is not empty and its length is within specified range.
func RequirePassword(pw string) *Reason {
	if r := RequireNotEmpty(pw, "password"); r != nil {
		return r
	}

	return RequireLenRange(pw, 8, 256, "password")
}
