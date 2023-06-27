package sandbox

import (
	"strings"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
)

type SignUpParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *SignUpParams) Validate() *render.ValidationError {
	p.Email = strings.TrimSpace(p.Email)
	p.Password = strings.TrimSpace(p.Password)

	ve := validator.New("email").Required().Email().Validate(p.Email)
	if ve != nil {
		return ve
	}

	ve = validator.New("password").Required().Validate(p.Password)
	if ve != nil {
		return ve
	}

	if !strings.HasSuffix(p.Email, ".test@ftchinese.com") {
		return &render.ValidationError{
			Message: "Only email addressing ending with .sandbox@ftchinese.com is allowed.",
			Field:   "email",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

type PasswordParams struct {
	Password string `json:"password"`
}

func (a *PasswordParams) Validate() *render.ValidationError {
	a.Password = strings.TrimSpace(a.Password)

	return validator.
		New("password").
		Required().
		Validate(a.Password)
}

type TestAccount struct {
	FtcID         string `json:"id" db:"ftc_id" gorm:"primaryKey"`
	Email         string `json:"email" db:"email"`
	ClearPassword string `json:"password" db:"clear_password"`
	CreatedBy     string `json:"createdBy" db:"created_by"`
}

func (a TestAccount) TableName() string {
	return "user_db.sandbox_account"
}

func (a TestAccount) WithPassword(pw string) TestAccount {
	a.ClearPassword = pw
	return a
}

type BaseAccount struct {
	FtcID string `json:"id"`
}

func (a BaseAccount) NewTestAccount(p SignUpParams, creator string) TestAccount {
	return TestAccount{
		FtcID:         a.FtcID,
		Email:         p.Email,
		ClearPassword: p.Password,
		CreatedBy:     creator,
	}
}
