package staff

import (
	"strings"
)

// Login specifies the the fields used for authentication
type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// Sanitize removes leading and trailing space of each field
func (l *Login) Sanitize() {
	l.UserName = strings.TrimSpace(l.UserName)
	l.Password = strings.TrimSpace(l.Password)
}


