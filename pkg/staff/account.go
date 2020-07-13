package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/letter"
)

// Account contains essential data of a user.
// It is used as response data for user authentication.
// It is also used to create a new user. In this case, password is set to a random string and sent to the Email of this new user. You must make sure the email already works.
type Account struct {
	ID           null.String `json:"id" db:"staff_id"`
	UserName     string      `json:"userName" db:"user_name"`             // Required, unique. Used for login.
	Email        string      `json:"email" db:"email"`                    // Required, unique.
	DisplayName  null.String `json:"displayName" db:"display_name"`       // Optional, unique.
	Department   null.String `json:"department" db:"department"`          // Optional.
	GroupMembers int64       `json:"groupMembers" db:"group_memberships"` // Required.
	IsActive     bool        `json:"isActive" db:"is_active"`
}

// Update updates all fields of an account by admin.
func (a Account) Update(input InputData) Account {
	a.UserName = input.UserName
	a.Email = input.Email
	a.DisplayName = input.DisplayName
	a.Department = input.Department
	a.GroupMembers = input.GroupMembers

	return a
}

// NormalizedName pick a title to name user in email.
func (a Account) NormalizeName() string {
	if a.DisplayName.Valid {
		return a.DisplayName.String
	}

	return a.UserName
}

// PasswordResetParcel create an email to enable resetting password.
func (a Account) PasswordResetParcel(session PwResetSession) (postoffice.Parcel, error) {
	body, err := letter.RenderPasswordReset(letter.CtxPasswordReset{
		DisplayName: a.NormalizeName(),
		URL:         session.BuildURL(),
	})

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "report@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "[FT中文网]重置密码",
		Body:        body,
	}, nil
}
