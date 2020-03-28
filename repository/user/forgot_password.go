package user

import (
	"gitlab.com/ftchinese/superyard/models/staff"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

// Password reset token
const stmtInsertResetToken = `
INSERT INTO backyard.password_reset
SET token = UNHEX(:token),
	email = :email,
	created_utc = UTC_TIMESTAMP()`

// SavePwResetToken send a password reset token to a user's email
func (env Env) SavePwResetToken(pr staff.PasswordReset) error {
	_, err := env.DB.NamedExec(stmtInsertResetToken, pr)

	if err != nil {
		logger.WithField("trace", "Env.SavePwResetToken").Error(err)

		return err
	}

	return nil
}

const stmtAccountByResetToken = stmt.StaffAccount + `
FROM backyard.password_reset AS r
	JOIN backyard.staff AS s
	ON r.email = s.email
WHERE r.token = UNHEX(?)
  AND r.is_used = 0
  AND DATE_ADD(r.created_utc, INTERVAL r.expires_in SECOND) > UTC_TIMESTAMP()
  AND s.is_active = 1
LIMIT 1`

// AccountByResetToken finds an account by
// a password reset token.
// This is used when the user clicked the
// link contained in password reset email.
func (env Env) AccountByResetToken(token string) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, stmtAccountByResetToken, token)

	if err != nil {
		logger.WithField("trace", "Env.AccountByResetToken").Error(err)

		return staff.Account{}, err
	}

	return a, err
}

const stmtDisableResetToken = `
UPDATE backyard.password_reset
SET is_used = 1
WHERE token = UNHEX(?)
LIMIT 1`

// DeleteResetToken deletes a password reset token after it was used.
func (env Env) DisableResetToken(token string) error {
	_, err := env.DB.Exec(stmtDisableResetToken, token)
	if err != nil {
		logger.WithField("trace", "Env.DeleteResetToken").Error(err)

		return err
	}

	return nil
}
