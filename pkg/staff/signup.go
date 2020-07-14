package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/letter"
)

func GenStaffID() string {
	return "stf_" + rand.String(12)
}

type SignUp struct {
	Account
	Password string `json:"password" db:"password"`
}

func NewSignUp(input InputData) SignUp {
	input.ID = null.StringFrom(GenStaffID())

	return input.SignUp
}

func (s SignUp) SignUpParcel(sourceURL string) (postoffice.Parcel, error) {

	if sourceURL == "" {
		sourceURL = "https://superyard.ftchinese.com"
	}

	body, err := letter.RenderSignUp(letter.CtxSignUp{
		DisplayName: s.NormalizeName(),
		LoginName:   s.UserName,
		Password:    s.Password,
		LoginURL:    sourceURL,
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
