package staff

import (
	"fmt"
	"github.com/FTChinese/go-rest/postoffice"
	"strings"
	"text/template"

	"github.com/FTChinese/go-rest/view"
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
func (a *Account) Validate() *view.Reason {
	// Is email is missing, not valid email address, or exceed 80 chars?
	if r := util.RequireEmail(a.Email); r != nil {
		return r
	}

	// Is the length displayName is within 20?
	if r := util.RequireNotEmptyWithMax(a.UserName, 255, "userName"); r != nil {
		return r
	}

	// Is userName exists and is within 20 chars?
	return util.OptionalMaxLen(a.DisplayName, 255, "displayName")
}

// TokenHolder generates a token for a user.
func (a Account) TokenHolder() (TokenHolder, error) {
	return NewTokenHolder(a.Email)
}

func (a Account) SignupParcel(pw string) (postoffice.Parcel, error) {
	tmpl, err := template.New("verification").Parse(SignupLetter)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	data := struct {
		Account
		Password string
	}{
		a,
		pw,
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
		ToName:      a.DisplayName,
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
		ToName:      a.DisplayName,
		Subject:     "[FT中文网]重置密码",
		Body:        body.String(),
	}, nil
}

//func (a Account) sendResetToken(token string, endpoint string) error {
//	request := gorequest.New()
//
//	_, _, errs := request.Post(endpoint).
//		Send(map[string]string{
//			"userName": a.UserName,
//			"token":    token,
//			"email":    a.Email,
//		}).
//		End()
//
//	if errs != nil {
//		logger.WithField("location", "Send password reset letter").Error(errs)
//
//		return errs[0]
//	}
//
//	return nil
//}
//
//// SendPassword send password to user's email address upon creation
//func (a Account) SendPassword(pass string, endpoint string) error {
//	request := gorequest.New()
//
//	_, _, errs := request.Post(endpoint).
//		Send(map[string]string{
//			"userName": a.UserName,
//			"email":    a.Email,
//			"password": pass,
//		}).
//		End()
//
//	if errs != nil {
//		logger.WithField("location", "Send welcome letter to new staff").Error(errs)
//
//		return errs[0]
//	}
//
//	return nil
//}

// FindAccount gets an account by user name.
// Use `activeOnly` to limit active staff only or all.
func (env Env) FindAccount(userName string, activeOnly bool) (Account, error) {
	var activeStmt string
	if activeOnly {
		activeStmt = "AND is_active = 1"
	}
	query := fmt.Sprintf(`
	SELECT id AS id,
		username AS userName,
		IFNULL(email, '') AS email,
		IFNULL(display_name, '') AS displayName,
		IFNULL(department, '') AS department,
		group_memberships AS groups
	FROM backyard.staff
	WHERE username = ?
		%s	
	LIMIT 1`, activeStmt)

	var a Account
	err := env.DB.QueryRow(query, userName).Scan(
		&a.ID,
		&a.UserName,
		&a.Email,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
	)

	if err != nil {
		logger.WithField("location", "Staff authentication").Error(err)

		return a, err
	}

	return a, nil
}
