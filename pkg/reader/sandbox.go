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

type SandboxUser struct {
	FtcAccount
	Password   string      `json:"password,omitempty" db:"password"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

// NewSandboxUser creates a new ftc account based on sandbox input.
func NewSandboxUser(input SandboxInput, creator string) SandboxUser {
	return SandboxUser{
		FtcAccount: FtcAccount{
			IDs: IDs{
				FtcID:   null.StringFrom(uuid.New().String()),
				UnionID: null.String{},
			},
			StripeID: null.String{},
			Email:    null.StringFrom(input.Email),
			UserName: null.String{},
		},
		Password:   input.Password,
		CreatedBy:  creator,
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
	}
}

// SandboxAccount contains a sandbox user info and membership.
type SandboxAccount struct {
	SandboxUser
	Membership Membership `json:"membership"`
}
