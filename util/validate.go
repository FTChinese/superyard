package util

import (
	"unicode/utf8"

	validate "github.com/asaskevich/govalidator"
)

// IsEmpty tests if str length is zero
func IsEmpty(str string) bool {
	return str == ""
}

// IsLength tests if a string's length is within a range.
func IsLength(str string, min, max int) bool {
	if min > max {
		min, max = max, min
	}
	strLength := utf8.RuneCountInString(str)
	return strLength >= min && strLength <= max
}

// MinLength tests if a string's length is longer than min
func MinLength(str string, min int) bool {
	strLength := utf8.RuneCountInString(str)
	return strLength >= min
}

// MaxLength tests if a string's length is under max
func MaxLength(str string, max int) bool {
	strLength := utf8.RuneCountInString(str)
	return strLength <= max
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if IsEmpty(email) {
		return UnprocessableError{
			Field: "email",
			Code:  CodeMissingField,
		}
	}

	if !validate.IsEmail(email) {
		return UnprocessableError{
			Field: "email",
			Code:  CodeInvalid,
		}
	}

	return nil
}
