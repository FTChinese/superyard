package employee

import (
	"github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/models/validator"
	"strings"
)

// TokenHolder holds a unique token for an email address.
// The token is readonly once generated.
type TokenHolder struct {
	Email string `json:"email" db:"email"`
	Token string `json:"-" db:"token"`
}

func (t *TokenHolder) GenerateToken() error {
	token, err := gorest.RandomHex(32)
	if err != nil {
		return err
	}

	t.Token = token

	return nil
}

func (t *TokenHolder) Sanitize() {
	t.Email = strings.TrimSpace(t.Email)
}

func (t TokenHolder) Validate() *validator.InputError {
	return validator.New("email").Required().Email().Validate(t.Email)
}

// PasswordReset is used as marshal target when user tries to reset password via email
type PasswordReset struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

// Sanitize removes leading and trailing space of each field
func (r *PasswordReset) Sanitize() {
	r.Token = strings.TrimSpace(r.Token)
	r.Password = strings.TrimSpace(r.Password)
}

func (r PasswordReset) Validate() *validator.InputError {
	return validator.New("password").
		Required().
		Min(8).
		Max(256).
		Validate(r.Password)
}
