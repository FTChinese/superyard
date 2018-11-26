package staff

import (
	"strings"

	"gitlab.com/ftchinese/backyard-api/util"
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
func (p *Password) Validate() *util.Reason {
	if r := util.RequirePassword(p.Old); r != nil {
		return r
	}

	if r := util.RequirePassword(p.New); r != nil {
		return r
	}

	return nil
}
