package employee

import (
	"strings"
)

// Login specifies the the fields used for authentication
type Login struct {
	UserName string `json:"userName" db:"user_name"`
	Password string `json:"password" db:"password"`
}

// Sanitize removes leading and trailing space of each field
func (l *Login) Sanitize() {
	l.UserName = strings.TrimSpace(l.UserName)
	l.Password = strings.TrimSpace(l.Password)
}
