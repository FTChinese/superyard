package staff

import "strings"

// MyftAccount is the ftc account owned by a staff
type MyftAccount struct {
	ID    string `json:"myftId"`
	Email string `json:"myftEmail"`
	IsVIP bool   `json:"isVip"`
}

// MyftCredential contains data to login to FTC
type MyftCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Sanitize removes leading and trailing spaces
func (c *MyftCredential) Sanitize() {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)
}


