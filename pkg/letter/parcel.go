package letter

import (
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/superyard/pkg/reader"
	"gitlab.com/ftchinese/superyard/pkg/staff"
	"gitlab.com/ftchinese/superyard/pkg/subs"
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
		Subject:     "[Superyar]Reset Password",
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

func OrderConfirmedParcel(a reader.FtcAccount, result subs.ConfirmationResult) (postoffice.Parcel, error) {
	body, err := RenderOrderConfirmed(CtxConfirmOrder{
		Name:           a.NormalizedName(),
		OrderCreatedAt: result.Order.CreatedAt.StringCN(),
		OrderID:        result.Order.ID,
		OrderAmount:    result.Order.ReadableAmount(),
		PayMethod:      result.Order.PaymentMethod.StringCN(),
		OrderStartDate: result.Order.StartDate.String(),
		OrderEndDate:   result.Order.EndDate.String(),
		Tier:           result.Membership.Tier.StringCN(),
		ExpirationDate: result.Membership.ExpireDate.String(),
	})

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return postoffice.Parcel{
		FromAddress: "no-reply@ftchinese.com",
		FromName:    "FT中文网",
		ToAddress:   a.Email.String,
		ToName:      a.NormalizedName(),
		Subject:     "订阅订单已确认",
		Body:        body,
	}, nil
}