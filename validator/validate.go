package validator

import (
	"unicode/utf8"

	validate "github.com/asaskevich/govalidator"

	"gitlab.com/ftchinese/backyard-api/view"
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

// Email validates an email address
func Email(email string) error {
	if IsEmpty(email) {
		return view.UnprocessableError{
			Field: "email",
			Code:  view.CodeMissingField,
		}
	}

	if !validate.IsEmail(email) {
		return view.UnprocessableError{
			Field: "email",
			Code:  view.CodeInvalid,
		}
	}

	return nil
}
