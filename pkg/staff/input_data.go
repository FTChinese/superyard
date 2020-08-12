package staff

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
)

// InputDate is used to parse various form
// data for a staff.
// Login: userName + password
// Reset password letter: email + sourceUrl
// Reset password: token + password
// Update password: oldPassword + password
// Set email: email
// Sign up: userName + email + displayName? + department? + groupMembers + password + sourceUrl
type InputData struct {
	SignUp
	OldPassword string `json:"oldPassword"`
	Token       string `json:"token"`
	SourceURL   string `json:"sourceUrl"` // Login page, or password reset page.
}

func (i *InputData) ValidateUserName() *render.ValidationError {
	i.UserName = strings.TrimSpace(i.UserName)

	return validator.
		New("userName").
		Required().
		MaxLen(64).
		Validate(i.UserName)
}

// ValidatePassword ensures the password field is valid.
// `minLen` indicates whether minimum length is required
// on password field.
// It is not required when logging in for backward compatibility reason.
func (i *InputData) ValidatePassword(minLen bool) *render.ValidationError {
	i.Password = strings.TrimSpace(i.Password)

	v := validator.
		New("password").
		Required().
		MaxLen(64)

	if minLen {
		v = v.MinLen(8)
	}

	return v.Validate(i.Password)
}

// ValidateLogin requires userName + password fields.
func (i *InputData) ValidateLogin() *render.ValidationError {

	ie := i.ValidateUserName()
	if ie != nil {
		return ie
	}

	return i.ValidatePassword(false)
}

func (i *InputData) Login() Credentials {
	return Credentials{
		UserName: i.UserName,
		Password: i.Password,
	}
}

// ValidateEmail validates email field when user asks for a
// password reset letter, or setting the email field.
func (i *InputData) ValidateEmail() *render.ValidationError {
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

// ValidationPasswordReset validates token + password fields
// when user submitted request to reset password.
func (i *InputData) ValidatePasswordReset() *render.ValidationError {
	i.Token = strings.TrimSpace(i.Token)

	ve := validator.
		New("token").
		Required().
		Validate(i.Token)

	if ve != nil {
		return ve
	}

	return i.ValidatePassword(true)
}

// ValidatePwUpdater validates oldPassword + password
// fields upon changing password.
func (i *InputData) ValidatePwUpdater() *render.ValidationError {
	i.Password = strings.TrimSpace(i.Password)
	i.OldPassword = strings.TrimSpace(i.OldPassword)

	ve := validator.
		New("oldPassword").
		Required().
		Validate(i.OldPassword)

	if ve != nil {
		return ve
	}

	return i.ValidatePassword(true)
}

// ValidateDisplayName validates displayName field.
func (i *InputData) ValidateDisplayName() *render.ValidationError {
	n := strings.TrimSpace(i.DisplayName.String)
	i.DisplayName = null.NewString(n, n != "")

	return validator.
		New("displayName").
		MaxLen(64).
		Validate(i.DisplayName.String)
}

// ValidateAccounts validates fields updated by admin.
// Fields:
// userName + email + displayName? + department? + groupMembers
func (i *InputData) ValidateAccount() *render.ValidationError {
	if ve := i.ValidateUserName(); ve != nil {
		return ve
	}

	if ve := i.ValidateEmail(); ve != nil {
		return ve
	}

	if i.DisplayName.Valid {
		if ve := i.ValidateDisplayName(); ve != nil {
			return ve
		}
	}

	if i.GroupMembers == 0 {
		return &render.ValidationError{
			Message: "Group access rights should not be 0",
			Field:   "groupMembers",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

// ValidateSignUp validates fields to create a new account.
// Fields:
// userName + email + displayName? + department? + groupMembers + password
func (i *InputData) ValidateSignUp() *render.ValidationError {
	if ve := i.ValidateAccount(); ve != nil {
		return ve
	}

	return i.ValidatePassword(true)
}
