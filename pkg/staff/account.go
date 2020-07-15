package staff

import (
	"github.com/guregu/null"
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
