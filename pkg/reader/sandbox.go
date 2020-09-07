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

	if !strings.HasSuffix(i.Email, ".sandbox@ftchinese.com") {
		return &render.ValidationError{
			Message: "Only email addressing ending with .sandbox@ftchinese.com is allowed.",
			Field:   "email",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

// SandboxPasswordSchema is used to update password.
type SandboxPasswordSchema struct {
	FtcID    string `json:"-" db:"ftc_id"`
	Password string `json:"password" db:"password"`
}

func (s *SandboxPasswordSchema) Validate() *render.ValidationError {
	s.Password = strings.TrimSpace(s.Password)

	return validator.New("password").Required().Validate(s.Password)
}

// NewSandboxFtcAccount creates a new ftc account based on sandbox input.
func NewSandboxFtcAccount(input SandboxInput, creator string) FtcAccount {
	return FtcAccount{
		IDs: IDs{
			FtcID:   null.StringFrom(uuid.New().String()),
			UnionID: null.String{},
		},
		StripeID:   null.String{},
		Email:      null.StringFrom(input.Email),
		UserName:   null.String{},
		Password:   input.Password,
		CreatedBy:  creator,
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
	}
}
