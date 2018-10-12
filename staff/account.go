package staff

import (
	"strings"

	"github.com/parnurzeal/gorequest"
	"gitlab.com/ftchinese/backyard-api/util"
)

// Role(s) a staff can have
const (
	RoleRoot      = 1
	RoleDeveloper = 2
	RoleEditor    = 4
	RoleWheel     = 8
	RoleSales     = 16
	RoleMarketing = 32
	RoleMetting   = 64
)

// Account contains essential data of a user.
// It is used as response data for user authenticztion.
// It is also used to create a new user. In this case, password is set to a random string and sent to the Email of this new user. You must make sure the email already works.
type Account struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`        // Required, unique, max 255 chars.
	UserName     string `json:"userName"`     // Required, unique, max 255 chars. Used for login.
	DisplayName  string `json:"displayName"`  // Optional, unique max 255 chars.
	Department   string `json:"department"`   // Optional, max 255 chars.
	GroupMembers int    `json:"groupMembers"` // Required.
}

// Sanitize removes leading and trailing spaces
func (a *Account) Sanitize() {
	a.Email = strings.TrimSpace(a.Email)
	a.UserName = strings.TrimSpace(a.UserName)
	a.DisplayName = strings.TrimSpace(a.DisplayName)
	a.Department = strings.TrimSpace(a.Department)
}

// Validate checks if required fields are valid
func (a Account) Validate() util.ValidationResult {
	// Is email is missing, not valid email address, or exceed 80 chars?
	if r := util.ValidateEmail(a.Email); r.IsInvalid {
		return r
	}

	if r := util.ValidateIsEmpty(a.UserName, "userName"); r.IsInvalid {
		return r
	}
	// Is the length displayName is within 20?
	if r := util.ValidateMaxLen(a.UserName, 255, "userName"); r.IsInvalid {
		return r
	}

	// Is userName exists and is within 20 chars?
	return util.ValidateMaxLen(a.DisplayName, 255, "displayName")
}

func (a Account) sendResetToken(token string, endpoint string) error {
	request := gorequest.New()

	_, _, errs := request.Post(endpoint).
		Send(map[string]string{
			"userName": a.UserName,
			"token":    token,
			"email":    a.Email,
		}).
		End()

	if errs != nil {
		logger.WithField("location", "Send password reset letter").Error(errs)

		return errs[0]
	}

	return nil
}

// SendPassword send password to user's email address upon creation
func (a Account) SendPassword(pass string, endpoint string) error {
	request := gorequest.New()

	_, _, errs := request.Post(endpoint).
		Send(map[string]string{
			"userName": a.UserName,
			"email":    a.Email,
			"password": pass,
		}).
		End()

	if errs != nil {
		logger.WithField("location", "Send welcome letter to new staff").Error(errs)

		return errs[0]
	}

	return nil
}
