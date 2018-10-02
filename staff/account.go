package staff

import "github.com/parnurzeal/gorequest"

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
	Email        string `json:"email"`
	UserName     string `json:"userName"`
	DisplayName  string `json:"displayName"`
	Department   string `json:"department"`
	GroupMembers int    `json:"groupMembers"`
}

func (a Account) sendResetToken(token string, endpoint string) error {
	request := gorequest.New()

	_, _, errs := request.Post(endpoint).
		Send(map[string]string{
			"userName": a.UserName,
			"token":    token,
			"address":  a.Email,
		}).
		End()

	if errs != nil {
		staffLogger.WithField("location", "Send password reset letter").Error(errs)

		return errs[0]
	}

	return nil
}
