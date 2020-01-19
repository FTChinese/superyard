package employee

import (
	"github.com/FTChinese/go-rest/view"
	"strings"
)

// Login specifies the the fields used for authentication
type Login struct {
	UserName string `json:"userName" db:"user_name"`
	Password string `json:"password" db:"password"`
}

func (l *Login) Validate() *view.Reason {
	if l.UserName == "" {
		return &view.Reason{
			Message: "User name is missing",
			Field:   "userName",
			Code:    view.CodeMissingField,
		}
	}

	if l.Password == "" {
		return &view.Reason{
			Message: "Password is missing",
			Field:   "password",
			Code:    view.CodeMissingField,
		}
	}

	return nil
}

// Sanitize removes leading and trailing space of each field
func (l *Login) Sanitize() {
	l.UserName = strings.TrimSpace(l.UserName)
	l.Password = strings.TrimSpace(l.Password)
}
