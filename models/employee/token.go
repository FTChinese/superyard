package employee

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/models/util"
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

func (t TokenHolder) Validate() *view.Reason {
	return util.RequireEmail(t.Email)
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

func (r PasswordReset) Validate() *view.Reason {
	if reason := util.RequireNotEmpty(r.Password, "password"); reason != nil {
		return reason
	}

	return util.RequirePassword(r.Password)
}
