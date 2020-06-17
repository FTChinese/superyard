package staff

import (
	"github.com/FTChinese/go-rest/render"
	"gitlab.com/ftchinese/superyard/pkg/validator"
	"strings"
)

// Login specifies the the fields used for authentication
type Login struct {
	UserName string `json:"userName" db:"user_name"`
	Password string `json:"password" db:"password"`
}

func (l *Login) Validate() *render.ValidationError {
	ie := validator.New("userName").Required().Validate(l.UserName)
	if ie != nil {
		return ie
	}

	return validator.New("password").Required().Validate(l.Password)
}

// Sanitize removes leading and trailing space of each field
func (l *Login) Sanitize() {
	l.UserName = strings.TrimSpace(l.UserName)
	l.Password = strings.TrimSpace(l.Password)
}
