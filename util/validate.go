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
	if isEmpty(value) {
		return ValidationResult{
			Field:     field,
			Code:      CodeMissingField,
			IsInvalid: true,
		}
	}

	return ValidationResult{}
}

// ValidateEmail makes sure an email is not empty, and is a valid email address
func ValidateEmail(email string) ValidationResult {

	if result := ValidateIsEmpty(email, "email"); result.IsInvalid {
		return result
	}

	if !validate.IsEmail(email) {
		return ValidationResult{
			Field:     "email",
			Code:      CodeInvalid,
			IsInvalid: true,
		}
	}

	return ValidationResult{}
}

// ValidatePassword makes sure the password is not empty, not execeeding max length
func ValidatePassword(pass string) ValidationResult {

	if result := ValidateIsEmpty(pass, "password"); result.IsInvalid {
		return result
	}

	if result := ValidateMaxLen(pass, 128, "password"); result.IsInvalid {
		return result
	}

	return ValidationResult{}
}
