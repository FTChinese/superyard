package validator

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"log"
)

const (
	msgTooLong  = "The length of %s should not exceed %d chars"
	msgTooShort = "The length of %s should not less than %d chars"
	msgLenRange = "The length of %s must be within %d to %d chars"
)

type Validator struct {
	fieldName  string
	isRequired bool
	min        int
	max        int
	isEmail    bool
	isURL      bool
}

func New(name string) *Validator {
	return &Validator{
		fieldName: name,
	}
}

func (v *Validator) Required() *Validator {
	v.isRequired = true
	return v
}

func (v *Validator) Min(min int) *Validator {
	v.min = min
	return v
}

func (v *Validator) Max(max int) *Validator {
	v.max = max
	return v
}

func (v *Validator) Range(min, max int) *Validator {
	v.min = min
	v.max = max
	return v
}

func (v *Validator) Email() *Validator {
	v.isEmail = true
	return v
}

func (v *Validator) URL() *Validator {
	v.isURL = true
	return v
}

func (v *Validator) Validate(value string) *InputError {
	if v.isEmail && v.isURL {
		log.Fatal("The validated value cannot be both an email and url")
	}

	if v.isRequired && !Required(value) {
		return &InputError{
			Message: "Missing required field",
			Field:   v.fieldName,
			Code:    CodeMissingField,
		}
	}

	if v.min > 0 && v.max > 0 && !StringInLength(value, v.min, v.max) {
		return &InputError{
			Message: fmt.Sprintf(msgLenRange, v.fieldName, v.min, v.max),
			Field:   v.fieldName,
			Code:    CodeInvalid,
		}
	}

	if v.min > 0 && !MinStringLength(value, v.min) {
		return &InputError{
			Message: fmt.Sprintf(msgTooShort, v.fieldName, v.min),
			Field:   v.fieldName,
			Code:    CodeInvalid,
		}
	}

	if v.max > 0 && !MaxStringLength(value, v.max) {
		return &InputError{
			Message: fmt.Sprintf(msgTooLong, v.fieldName, v.max),
			Field:   v.fieldName,
			Code:    CodeInvalid,
		}
	}

	if v.isEmail && !govalidator.IsEmail(value) {
		return &InputError{
			Message: "Invalid email address",
			Field:   v.fieldName,
			Code:    CodeInvalid,
		}
	}

	if v.isURL && !govalidator.IsURL(value) {
		return &InputError{
			Message: "Invalid URL",
			Field:   v.fieldName,
			Code:    CodeInvalid,
		}
	}

	return nil
}
