package staff

import (
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/models/employee"
)

func (env Env) RetrieveProfile(id string) (employee.Profile, error) {
	var p employee.Profile

	err := env.DB.Get(&p, stmtSelectProfile, id)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveProfile").Error(err)

		return p, err
	}

	return p, nil
}

func (env Env) ListStaff(p builder.Pagination) ([]employee.Profile, error) {
	profiles := make([]employee.Profile, 0)

	err := env.DB.Select(&profiles,
		stmtListStaff,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListStaff").Error(err)

		return profiles, err
	}

	return profiles, nil
}

// UpdateAccount updates a staff's account.
//
//	PATCH /admin/accounts/{name}
//
// Input {userName: string, email: string, displayName: string, department: string, groupMembers: number}
func (env Env) UpdateProfile(p employee.Profile) error {
	_, err := env.DB.NamedExec(stmtUpdateProfile, &p)
	if err != nil {
		logger.WithField("trace", "Env.UpdateAccount").Error(err)
		return err
	}

	return nil
}

// VerifyPassword verifies a staff's password and returns
// account data if it is correct.
func (env Env) VerifyPassword(a employee.Account) (employee.Account, error) {
	var account employee.Account
	err := env.DB.Get(&account, stmtVerifyPassword, a.ID, a.Password)

	if err != nil {
		logger.WithField("trace", "VerifyPassword").Error(err)
		return employee.Account{}, err
	}

	return account, nil
}

// Change password is used by both UpdatePassword after user logged in, or reset password if user forgot it.
func (env Env) changePassword(password string, userName string) error {
	tx, err := env.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(stmtUpdatePassword, password, userName)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.changePassword").Error(err)

		return err
	}

	_, err = tx.Exec(stmtUpdateLegacyPassword, password, userName)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.changePassword").Error(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.WithField("trace", "changePassword").Error(err)
		return err
	}

	return nil
}

// UpdatePassword allows user to change password in its settings.
func (env Env) UpdatePassword(a employee.Account) error {

	tx, err := env.DB.Beginx()
	if err != nil {
		return err
	}

	// Update password in the new table.
	_, err = tx.NamedExec(stmtUpdatePassword, a)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.UpdatePassword").Error(err)
		return err
	}

	// Update password in old table
	_, err = tx.NamedExec(stmtUpdateLegacyPassword, a)
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
