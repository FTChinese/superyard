package sandbox

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"strings"
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
	FtcID         string `json:"id" db:"ftc_id"`
	Email         string `json:"email" db:"email"`
	ClearPassword string `json:"password" db:"clear_password"`
	CreatedBy     string `json:"createdBy" db:"created_by"`
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

// StmtInsertTestAccount records which account is sandbox and store the password as clear text.
const StmtInsertTestAccount = `
INSERT INTO user_db.sandbox_account
SET ftc_id = :ftc_id,
	email = :email,
	clear_password = :clear_password,
	created_by = :created_by
`

const colTestAccount = `
SELECT ftc_id,
	email,
	clear_password,
	created_by
FROM user_db.sandbox_account
`

const StmtRetrieveTestUser = colTestAccount + `
WHERE ftc_id = ?
LIMIT 1
`

// StmtListTestUsers retrieves a list of FtcAccount.
const StmtListTestUsers = colTestAccount + `
ORDER BY email
LIMIT ? OFFSET ?
`

const StmtCountTestUser = `
SELECT COUNT(*) AS row_count
FROM user_db.sandbox_account
`

const StmtUpdateTestUserPassword = `
UPDATE user_db.sandbox_account
SET clear_password = :clear_password
WHERE ftc_id = :ftc_id
LIMIT 1`

const StmtUpdatePassword = `
UPDATE cmstmp01.userinfo
SET password := MD5(:clear_password),
	updated_utc := UTC_TIMESTAMP()
WHERE user_id = :ftc_id
LIMIT 1
`

const StmtDeleteTestUser = `
DELETE FROM user_db.sandbox_account
WHERE ftc_id = ?
LIMIT 1
`
