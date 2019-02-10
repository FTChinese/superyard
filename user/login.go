package user

import "strings"

// Login contains data to login to FTC
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Sanitize removes leading and trailing spaces
func (c *Login) Sanitize() {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)
}