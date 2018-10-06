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
func ValidateLength(value string, min int, max int, field string) ValidationResult {
	if min > 0 && value == "" {
		return ValidationResult{
			Field:     field,
			Code:      CodeMissingField,
			IsInvalid: true,
		}
	}

	if !isLength(value, min, max) {
		return ValidationResult{
			Message:   fmt.Sprintf("The length of %s should be within %d to %d chars", value, min, max),
			Field:     field,
			Code:      CodeInvalid,
			IsInvalid: true,
		}
	}

	return ValidationResult{}
}

// ValidateMaxLen makes sure the value's length does not exceed the max limit
func ValidateMaxLen(value string, max int, field string) ValidationResult {
	if !maxLength(value, max) {
		return ValidationResult{
			Message:   fmt.Sprintf("The length of %s should not exceed %d chars", field, max),
			Field:     field,
			Code:      CodeInvalid,
			IsInvalid: true,
		}
	}

	return ValidationResult{}
}

// ValidateIsEmpty makes sure the value is not an empty string
func ValidateIsEmpty(value string, field string) ValidationResult {
	if value == "" {
		return ValidationResult{
			Field:     field,
			Code:      CodeMissingField,
			IsInvalid: true,
		}
	}

	return ValidationResult{}
}

// ValidateEmail makes sure an email is a valid email address, and max length does not exceed 20 chars
func ValidateEmail(email string) ValidationResult {

	if r := ValidateIsEmpty(email, "email"); r.IsInvalid {
		return r
	}

	if !validate.IsEmail(email) {
		return ValidationResult{
			Field:     "email",
			Code:      CodeInvalid,
			IsInvalid: true,
		}
	}

	return ValidateMaxLen(email, 20, "email")
}

// ValidatePassword makes sure the password is not empty, and not execeeds max length
func ValidatePassword(pass string) ValidationResult {
	return ValidateLength(pass, 8, 128, "password")
}
