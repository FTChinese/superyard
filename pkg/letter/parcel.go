package letter

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/staff"
)

func PasswordResetParcel(a staff.Account, session staff.PwResetSession) (postoffice.Parcel, error) {
	body, err := RenderPasswordReset(CtxPasswordReset{
		DisplayName: a.NormalizeName(),
		URL:         session.BuildURL(),
	})

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "report@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email,
		ToName:      a.NormalizeName(),
		Subject:     "[Superyard]Reset Password",
		Body:        body,
	}, nil
}

func SignUpParcel(s staff.SignUp, sourceURL string) (postoffice.Parcel, error) {
	if sourceURL == "" {
		sourceURL = "https://superyard.ftchinese.com"
	}

	body, err := RenderSignUp(CtxSignUp{
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

func MemberUpsertParcel(a reader.Account) (postoffice.Parcel, error) {
	body, err := RenderUpsertMember(CtxUpsertMember{
		Name:           a.NormalizedName(),
		Tier:           a.Membership.Tier.StringCN(),
		ExpirationDate: a.Membership.ExpireDate.String(),
	})

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email.String,
		ToName:      a.NormalizedName(),
		Subject:     "会员状态变动",
		Body:        body,
	}, nil
}
