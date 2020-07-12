package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"strings"
	"text/template"
)

// SignUp creates a new employee.
type SignUp struct {
	PasswordHolder
	BaseAccount
}

// NewSignUp creates a new user based on submitted data.
func NewSignUp(input InputData) SignUp {
	input.IsActive = true

	return SignUp{
		PasswordHolder: PasswordHolder{
			ID:       GenStaffID(),
			Password: input.Password,
		},
		BaseAccount: input.BaseAccount,
	}
}

func (s SignUp) SignUpParcel() (postoffice.Parcel, error) {
	tmpl, err := template.New("verification").Parse(SignupLetter)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	var body strings.Builder
	err = tmpl.Execute(&body, s)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "report@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   s.Email,
		ToName:      s.NormalizeName(),
		Subject:     "Welcome",
		Body:        body.String(),
	}, nil
}
