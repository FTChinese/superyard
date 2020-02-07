package employee

import (
	"gitlab.com/ftchinese/superyard/models/validator"
	"strings"
)

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
func (p *Password) Validate() *validator.InputError {
	ie := validator.New("oldPassword").Required().Min(1).Max(256).Validate(p.Old)
	if ie != nil {
		return ie
	}

	return validator.New("newPassword").Required().Min(8).Max(256).Validate(p.New)
}
