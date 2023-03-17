package user

import (
	"strings"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
)

// Login specifies the the fields used for authentication
type Credentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func (c *Credentials) Validate() *render.ValidationError {
	c.UserName = strings.TrimSpace(c.UserName)

	err := validator.
		New("userName").
		Required().
		MaxLen(64).
		Validate(c.UserName)

	if err != nil {
		return err
	}

	c.Password = strings.TrimSpace(c.Password)

	return validator.
		New("password").
		Required().
		MaxLen(64).
		Validate(c.Password)

}

type ParamsEmail struct {
	Email string `json:"email"`
}

// ValidateEmail validates email field when user asks for a
// password reset letter, or setting the email field.
func (i *ParamsEmail) Validate() *render.ValidationError {
	i.Email = strings.TrimSpace(i.Email)

	ve := validator.
		New("email").
		Required().
		Email().
		Validate(i.Email)

	if ve != nil {
		return ve
	}

	if !strings.HasSuffix(i.Email, "@ftchinese.com") {
		return &render.ValidationError{
			Message: "Email must be owned by ftchinese",
			Field:   "email",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

type ParamsDisplayName struct {
	DisplayName string `json:"displayName"`
}

// ValidateDisplayName validates displayName field.
func (i *ParamsDisplayName) Validate() *render.ValidationError {

	i.DisplayName = strings.TrimSpace(i.DisplayName)

	return validator.
		New("displayName").
		MaxLen(64).
		Validate(i.DisplayName)
}

type ParamsPasswords struct {
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

func (i *ParamsPasswords) Validate() *render.ValidationError {
	i.Password = strings.TrimSpace(i.Password)
	i.OldPassword = strings.TrimSpace(i.OldPassword)

	ve := validator.
		New("oldPassword").
		Required().
		Validate(i.OldPassword)

	if ve != nil {
		return ve
	}

	return validator.
		New("password").
		Required().
		MaxLen(64).
		MinLen(8).
		Validate(i.Password)
}

type ParamsForgotPassLetter struct {
	ParamsEmail
	SourceURL string `json:"sourceUrl"` // Optional. From which website the request is sent so that you can build verification url based on it
}

type ParamsResetPass struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Validation validates token + password fields
// when user submitted request to reset password.
func (i *ParamsResetPass) Validate() *render.ValidationError {
	i.Token = strings.TrimSpace(i.Token)

	ve := validator.
		New("token").
		Required().
		Validate(i.Token)

	if ve != nil {
		return ve
	}

	return validator.
		New("password").
		Required().
		MaxLen(64).
		MinLen(8).
		Validate(i.Password)
}
