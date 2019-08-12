package reader

import "strings"

// Login contains data to login to FTC.
// This is used when a CMS user attempts associated CMS account with FTC account.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Sanitize removes leading and trailing spaces
func (c *Login) Sanitize() {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)
}
