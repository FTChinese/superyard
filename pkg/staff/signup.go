package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/superyard/pkg/letter"
)

func SignUpParcel(s InputData) (postoffice.Parcel, error) {
	body, err := letter.RenderSignUp(letter.CtxSignUp{
		DisplayName: s.NormalizeName(),
		LoginName:   s.UserName,
		Password:    s.Password,
		LoginURL:    s.SourceURL,
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
