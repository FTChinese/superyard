package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"strings"
)

// SandboxInput is used to parse request body to create a sandbox account.
type SandboxInput struct {
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

	return nil
}

type SandboxAccount struct {
	FtcAccount
	Password   string      `json:"-" db:"password"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"createdUtc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updatedUtc"`
}

// NewSandboxAccount creates a new ftc account based on sandbox input.
func NewSandboxAccount(input SandboxInput) SandboxAccount {
	return SandboxAccount{
		FtcAccount: FtcAccount{
			FtcID:    null.StringFrom(uuid.New().String()),
			UnionID:  null.String{},
			StripeID: null.String{},
			Email:    null.StringFrom(input.Email),
			UserName: null.String{},
		},
		Password:   input.Password,
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
	}
}
