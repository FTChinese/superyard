package user

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

const stmtUpdatePassword = `
UPDATE backyard.staff
SET password = UNHEX(MD5(:password)),
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = :staff_id
	AND is_active = 1
LIMIT 1`

const stmtUpdateLegacyPassword = `
UPDATE cmstmp01.managers
	SET password = MD5(:password)
WHERE username = :user_name
LIMIT 1`

// UpdatePassword allows user to change password.
// It also updates the legacy table, which does
// not have a staff_id column. So we use user_name
// to update the legacy table.
// Therefore, to update password, we should know
// user'd id and user name.
func (env Env) UpdatePassword(c staff.Credentials) error {

	tx, err := env.DB.Beginx()
	if err != nil {
		return err
	}

	// Update password in the new table.
	_, err = tx.NamedExec(stmtUpdatePassword, c)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.UpdatePassword").Error(err)
		return err
	}

	// Update password in old table
	_, err = tx.NamedExec(stmtUpdateLegacyPassword, c)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.UpdatePassword").Error(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.WithField("trace", "Env.UpdatePassword").Error(err)
		return err
	}

	return nil
}

const stmtVerifyPassword = stmt.StaffAccount + `
FROM backyard.staff AS s
WHERE (s.staff_id, s.password) = (?, UNHEX(MD5(?)))
	AND s.is_active = 1`

// VerifyPassword verifies a staff's password
// when user tries to change password.
// ID and Password fields are required.
func (env Env) VerifyPassword(c staff.Credentials) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, stmtVerifyPassword, c.ID, c.Password)

	if err != nil {
		logger.WithField("trace", "VerifyPassword").Error(err)
		return staff.Account{}, err
	}

	return a, nil
}
