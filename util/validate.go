package util

import (
	"fmt"
	"unicode/utf8"

	validate "github.com/asaskevich/govalidator"
)

// isEmpty tests if str length is zero
func isEmpty(str string) bool {
	return str == ""
}

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

// ValidateLength makes sure the value's length is within the specified range
func ValidateLength(value string, min int, max int, field string) InvalidReason {

	if !isLength(value, min, max) {
		return InvalidReason{
			Message:   fmt.Sprintf("The length of %s should be within %d to %d chars", value, min, max),
			Field:     field,
			Code:      CodeInvalid,
			IsInvalid: true,
		}
	}

	return InvalidReason{}
}

// ValidateMaxLen makes sure the value's length does not exceed the max limit.
// Empty string is valid.
func ValidateMaxLen(value string, max int, field string) InvalidReason {
	if !maxLength(value, max) {
		return InvalidReason{
			Message:   fmt.Sprintf("The length of %s should not exceed %d chars", field, max),
			Field:     field,
			Code:      CodeInvalid,
			IsInvalid: true,
		}
	}

	return InvalidReason{}
}

// ValidateIsEmpty makes sure the value is not an empty string
func ValidateIsEmpty(value string, field string) InvalidReason {
	if value == "" {
		return InvalidReason{
			Field:     field,
			Code:      CodeMissingField,
			IsInvalid: true,
		}
	}

	return InvalidReason{}
}

// ValidateEmail makes sure an email is a valid email address, and max length does not exceed 80 chars
func ValidateEmail(email string) InvalidReason {

	if r := ValidateIsEmpty(email, "email"); r.IsInvalid {
		return r
	}

	if !validate.IsEmail(email) {
		return InvalidReason{
			Field:     "email",
			Code:      CodeInvalid,
			IsInvalid: true,
		}
	}

	return ValidateMaxLen(email, 254, "email")
}

// ValidatePassword makes sure the length of password is at least 8, and at most 255.
func ValidatePassword(pass string) InvalidReason {
	return ValidateLength(pass, 8, 255, "password")
}
