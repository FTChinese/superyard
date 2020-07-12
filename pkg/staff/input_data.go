package staff

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/validator"
	"strings"
)

// InputDate is used to parse various form
// data for a staff.
// Login: userName + password
// Reset password letter: email
// Reset password: token + password
// Update password: oldPassword + password
// Set email: email
// Sign up: userName + email + displayName? + department? + groupMembers + password
type InputData struct {
	BaseAccount
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
	Token       string `json:"token"`
}

func (i *InputData) ValidateUserName() *render.ValidationError {
	i.UserName = strings.TrimSpace(i.UserName)

	return validator.
		New("userName").
		Required().
		Validate(i.UserName)
}

// ValidateLogin requires userName + password fields.
func (i *InputData) ValidateLogin() *render.ValidationError {
	i.UserName = strings.TrimSpace(i.UserName)
	i.Password = strings.TrimSpace(i.Password)

	ie := i.ValidateUserName()
	if ie != nil {
		return ie
	}

	return validator.New("password").Required().Validate(i.Password)
}

func (i *InputData) Login() Login {
	return Login{
		UserName: i.UserName,
		Password: i.Password,
	}
}

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

// ValidatePwUpdater validates fields upon changing password.
// Previously (2020-06-04) client request body contains field:
// oldPassword + newPassword.
// Then we changed the request body to:
// oldPassword + password.
// To keep backward compatibility, we should manually copy
// newPassword to password if it exists.
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

func (i *InputData) ValidateDisplayName() *render.ValidationError {
	n := strings.TrimSpace(i.DisplayName.String)
	i.DisplayName = null.NewString(n, n != "")

	return validator.
		New("displayName").
		MaxLen(64).
		Validate(i.DisplayName.String)
}

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

func (i *InputData) ValidateSignUp() *render.ValidationError {
	if ve := i.ValidateAccount(); ve != nil {
		return ve
	}

	return i.ValidatePassword(true)
}
