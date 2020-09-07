package user

import (
	"github.com/FTChinese/superyard/pkg/staff"
)

// VerifyPassword verifies a staff's password
// when user tries to change password.
// ID and Password fields are required.
func (env Env) VerifyPassword(verifier staff.PasswordVerifier) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(
		&a,
		staff.StmtVerifyPassword,
		verifier.StaffID,
		verifier.OldPassword)

	if err != nil {
		return staff.Account{}, err
	}

	return a, nil
}

// UpdatePassword allows user to change password.
// It also updates the legacy table, which does
// not have a staff_id column. So we use user_name
// to update the legacy table.
// Therefore, to update password, we should know
// user'd id and user name.
func (env Env) UpdatePassword(holder staff.Credentials) error {

	tx, err := env.DB.Beginx()
	if err != nil {
		return err
	}

	// Update password in the new table.
	_, err = tx.NamedExec(staff.StmtUpdatePassword, holder)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Update password in old table
	_, err = tx.NamedExec(staff.StmtUpdateLegacyPassword, holder)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
