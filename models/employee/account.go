package employee

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/guregu/null"
	"strings"
	"text/template"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/models/util"
)

// Account contains essential data of a user.
// It is used as response data for user authentication.
// It is also used to create a new user. In this case, password is set to a random string and sent to the Email of this new user. You must make sure the email already works.
type Account struct {
	ID           string      `json:"id" db:"staff_id"`
	Email        string      `json:"email" db:"email"`        // Required, unique, max 255 chars.
	UserName     string      `json:"userName" db:"user_name"` // Required, unique, max 255 chars. Used for login.
	Password     string      `json:"-" db:"password"`
	IsActive     bool        `json:"isActive" db:"is_active"`
	DisplayName  null.String `json:"displayName" db:"display_name"`       // Optional, unique max 255 chars.
	Department   null.String `json:"department" db:"department"`          // Optional, max 255 chars.
	GroupMembers int64       `json:"groupMembers" db:"group_memberships"` // Required.
}

// NewAccount creates an account with password generated randomly.
func NewAccount() (Account, error) {
	password, err := gorest.RandomHex(4)
	if err != nil {
		return Account{}, err
	}

	id, err := gorest.RandomHex(8)
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:       "stf_" + id,
		Password: password,
	}, nil
}

func (a *Account) GenerateID() error {
	id, err := gorest.RandomHex(8)
	if err != nil {
		return err
	}

	a.ID = "stf_" + id

	return nil
}

func (a *Account) GeneratePassword() error {
	password, err := gorest.RandomHex(4)
	if err != nil {
		return err
	}

	a.Password = password

	return nil
}

func (a Account) NormalizeName() string {
	if a.DisplayName.Valid {
		return a.DisplayName.String
	}

	return a.UserName
}

// Sanitize removes leading and trailing spaces
func (a *Account) Sanitize() {
	a.Email = strings.TrimSpace(a.Email)
	a.UserName = strings.TrimSpace(a.UserName)
	//a.DisplayName = strings.TrimSpace(a.DisplayName)
	//a.Department = strings.TrimSpace(a.Department)
}

// Validate checks if required fields are valid
func (a Account) Validate() *view.Reason {
	// Is email is missing, not valid email address, or exceed 80 chars?
	if r := util.RequireEmail(a.Email); r != nil {
		return r
	}

	// Is the length displayName is within 20?
	if r := util.RequireNotEmptyWithMax(a.UserName, 255, "userName"); r != nil {
		return r
	}

	// Is userName exists and is within 20 chars?
	return util.OptionalMaxLen(a.DisplayName.String, 255, "displayName")
}

func (a Account) SignUpParcel() (postoffice.Parcel, error) {
	tmpl, err := template.New("verification").Parse(SignupLetter)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	var body strings.Builder
	err = tmpl.Execute(&body, a)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "report@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "Welcome",
		Body:        body.String(),
	}, nil
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
