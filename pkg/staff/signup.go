package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/superyard/pkg/letter"
)

// SignUp creates a new employee.
type SignUp struct {
	PasswordHolder
	BaseAccount
	LoginURL string
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
	body, err := letter.RenderSignUp(letter.CtxSignUp{
		DisplayName: s.NormalizeName(),
		LoginName:   s.UserName,
		Password:    s.Password,
		LoginURL:    s.LoginURL,
	})

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "report@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   s.Email,
		ToName:      s.NormalizeName(),
		Subject:     "Welcome",
		Body:        body,
	}, nil
}
