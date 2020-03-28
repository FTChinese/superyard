package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"strings"
	"text/template"
)

// SignUp creates a new employee.
type SignUp struct {
	Account
	Password string `db:"password"`
}

// NewSignUp creates a new user based on submitted data.
func NewSignUp(base BaseAccount) SignUp {
	return SignUp{
		Account: Account{
			ID:          null.StringFrom(GenStaffID()),
			BaseAccount: base,
			IsActive:    true,
		},
		Password: rand.String(8),
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
