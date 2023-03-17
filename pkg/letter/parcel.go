package letter

import (
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/postman"
)

func PasswordResetParcel(a user.Account, resetURL string) (postman.Parcel, error) {
	body, err := RenderPasswordReset(CtxPasswordReset{
		DisplayName: a.NormalizeName(),
		URL:         resetURL,
	})

	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		FromAddress: "report@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "[Superyard]Reset Password",
		Body:        body,
	}, nil
}
