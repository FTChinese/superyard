package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/validator"
	"strings"
	"text/template"
)

// BaseAccount contains the shared fields of account-related types.
type BaseAccount struct {
	UserName     string      `json:"userName" db:"user_name"`             // Required, unique, max 255 chars. Used for login.
	Email        string      `json:"email" db:"email"`                    // Required, unique, max 255 chars.
	DisplayName  null.String `json:"displayName" db:"display_name"`       // Optional, unique max 255 chars.
	Department   null.String `json:"department" db:"department"`          // Optional, max 255 chars.
	GroupMembers int64       `json:"groupMembers" db:"group_memberships"` // Required.
}

// Sanitize removes leading and trailing spaces
func (a *BaseAccount) Sanitize() {
	a.Email = strings.TrimSpace(a.Email)
	a.UserName = strings.TrimSpace(a.UserName)
	//a.DisplayName = strings.TrimSpace(a.DisplayName)
	//a.Department = strings.TrimSpace(a.Department)
}

func (a BaseAccount) ValidateEmail() *render.ValidationError {
	return validator.New("email").Required().Max(256).Email().Validate(a.Email)
}

func (a BaseAccount) ValidateDisplayName() *render.ValidationError {
	return validator.New("displayName").Max(256).Validate(a.DisplayName.String)
}

// Validate checks if required fields are valid
func (a BaseAccount) Validate() *render.ValidationError {
	ve := a.ValidateEmail()
	if ve != nil {
		return ve
	}

	ve = validator.New("userName").Required().Max(256).Validate(a.UserName)
	if ve != nil {
		return ve
	}

	return a.ValidateDisplayName()
}

// NormalizedName pick a title to name user in email.
func (a BaseAccount) NormalizeName() string {
	if a.DisplayName.Valid {
		return a.DisplayName.String
	}

	return a.UserName
}

// Account contains essential data of a user.
// It is used as response data for user authentication.
// It is also used to create a new user. In this case, password is set to a random string and sent to the Email of this new user. You must make sure the email already works.
type Account struct {
	ID null.String `json:"id" db:"staff_id"`
	BaseAccount
	IsActive bool `json:"isActive" db:"is_active"`
}

// Credentials build a Credentials instance for
// current account.
func (a Account) Credentials(pw string) Credentials {
	return Credentials{
		ID: a.ID.String,
		Login: Login{
			UserName: a.UserName,
			Password: pw,
		},
	}
}

// SetEmail sets the email field.
func (a *Account) SetEmail(email string) *render.ValidationError {
	if a.Email != "" {
		return &render.ValidationError{
			Message: "email could only be set once",
			Field:   "email",
			Code:    render.CodeAlreadyExists,
		}
	}

	a.Email = email

	return nil
}

func (a Account) PasswordResetParcel(token string) (postoffice.Parcel, error) {
	tmpl, err := template.New("verification").Parse(PasswordResetLetter)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	data := struct {
		Account
		Token string
	}{
		a,
		token,
	}
	var body strings.Builder
	err = tmpl.Execute(&body, data)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "report@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "[FT中文网]重置密码",
		Body:        body.String(),
	}, nil
}
