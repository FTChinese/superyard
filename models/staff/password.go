package staff

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"gitlab.com/ftchinese/superyard/models/validator"
	"strings"
)

// PasswordReset holds password resetting data.
// The fields won't all exist at the same time.
// When requesting a reset email, only `Email`
// is present.
// When the link in email is clicked, only `Token`
// is present.
// When new password is submitted, `Token` and `Password` will be present.
type PasswordReset struct {
	Email    string `json:"email" db:"email"`
	Token    string `json:"token" db:"token"`
	Password string `json:"password"`
}

// GeneratedToken creates a password reset token
// after user's email is submitted.
func (r *PasswordReset) GenerateToken() error {
	token, err := gorest.RandomHex(32)
	if err != nil {
		return err
	}

	r.Token = token

	return nil
}

// Sanitize removes leading and trailing space of each field
func (r *PasswordReset) Sanitize() {
	r.Token = strings.TrimSpace(r.Token)
	r.Password = strings.TrimSpace(r.Password)
	r.Token = strings.TrimSpace(r.Token)
}

func (r PasswordReset) ValidateEmail() *render.ValidationError {
	return validator.New("email").Required().Email().Validate(r.Email)
}

func (r PasswordReset) ValidatePass() *render.ValidationError {
	return validator.New("password").
		Required().
		Min(8).
		Max(256).
		Validate(r.Password)
}

// Credentials holds user's identity and password.
type Credentials struct {
	ID string `db:"staff_id"`
	Login
}

// Password marshals request data for updating password
type Password struct {
	Old string `json:"oldPassword"`
	New string `json:"newPassword"`
}

// Sanitize removes leading and  trailing white spaces.
func (p *Password) Sanitize() {
	p.Old = strings.TrimSpace(p.Old)
	p.New = strings.TrimSpace(p.New)
}

// Validate checks if old and new password are valid
func (p *Password) Validate() *render.ValidationError {
	ie := validator.New("oldPassword").Required().Min(1).Max(256).Validate(p.Old)
	if ie != nil {
		return ie
	}

	return validator.New("newPassword").Required().Min(8).Max(256).Validate(p.New)
}
