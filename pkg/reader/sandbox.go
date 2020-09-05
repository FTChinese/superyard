package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"strings"
)

// SandboxInput is used to parse request body to create a sandbox account.
type SandboxInput struct {
	FtcID    string `json:"ftcId"` // Only used when changing password.
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (i *SandboxInput) Validate() *render.ValidationError {
	i.Email = strings.TrimSpace(i.Email)
	i.Password = strings.TrimSpace(i.Password)

	ve := validator.New("email").Required().Email().Validate(i.Email)
	if ve != nil {
		return ve
	}

	ve = validator.New("password").Required().Validate(i.Password)
	if ve != nil {
		return ve
	}

	if !strings.HasSuffix(i.Email, ".sandbox@ftchinese.com") {
		return &render.ValidationError{
			Message: "Only email addressing ending with .sandbox@ftchinese.com is allowed.",
			Field:   "email",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

// SandboxFtcAccount contains the shared fields used for SandboxAccount
// and its flattened schema.
type SandboxFtcAccount struct {
	FtcAccount
	Password  string `json:"password" db:"password"`
	CreatedBy string `json:"createdBy" db:"created_by"`
}

// NewSandboxFtcAccount creates a new ftc account based on sandbox input.
func NewSandboxFtcAccount(input SandboxInput, creator string) SandboxFtcAccount {
	return SandboxFtcAccount{
		FtcAccount: FtcAccount{
			IDs: IDs{
				FtcID:   null.StringFrom(uuid.New().String()),
				UnionID: null.String{},
			},
			StripeID:   null.String{},
			Email:      null.StringFrom(input.Email),
			UserName:   null.String{},
			CreatedUTC: chrono.TimeNow(),
			UpdatedUTC: chrono.TimeNow(),
		},
		Password:  input.Password,
		CreatedBy: creator,
	}
}

// SandboxAccount contains a sandbox user info and membership.
type SandboxAccount struct {
	SandboxFtcAccount
	Kind       enum.AccountKind `json:"kind"`
	Wechat     Wechat           `json:"wechat"`
	Membership Membership       `json:"membership"`
}

type SandboxJoinedAccountSchema struct {
	SandboxFtcAccount
	Wechat
	VIP bool `db:"is_vip"`
}

func (s SandboxJoinedAccountSchema) Build(m Membership) SandboxAccount {
	if s.VIP {
		m.Tier = enum.TierVIP
	}

	return SandboxAccount{
		SandboxFtcAccount: s.SandboxFtcAccount,
		Wechat:            s.Wechat,
		Kind:              enum.AccountKindFtc,
		Membership:        m,
	}
}
