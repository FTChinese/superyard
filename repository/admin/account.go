package admin

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

const stmtAccountByID = stmt.StaffAccount + `
FROM backyard.staff AS s
WHERE s.staff_id = ?
LIMIT 1`

// AccountByID retrieves staff account by
// email column.
func (env Env) AccountByID(id string) (staff.Account, error) {
	var a staff.Account

	if err := env.DB.Get(&a, stmtAccountByID, id); err != nil {
		return staff.Account{}, err
	}

	return a, nil
}

const stmtAccountByName = stmt.StaffAccount + `
FROM backyard.staff AS s
WHERE s.user_name = ?
LIMIT 1`

// AccountByName loads an account when by name
// is submitted to request a password reset letter.
func (env Env) AccountByName(name string) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, stmtAccountByName, name)

	if err != nil {
		logger.WithField("trace", "Env.AccountByName").Error(err)

		return staff.Account{}, err
	}

	return a, err
}

const stmtListStaff = stmt.StaffAccount + `
FROM backyard.staff AS s
ORDER BY s.user_name ASC
LIMIT ? OFFSET ?`

func (env Env) ListStaff(p util.Pagination) ([]staff.Account, error) {
	accounts := make([]staff.Account, 0)

	err := env.DB.Select(&accounts,
		stmtListStaff,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListStaff").Error(err)

		return accounts, err
	}

	return accounts, nil
}

const stmtUpdateProfile = `
UPDATE backyard.staff
SET user_name = :user_name,
	email = :email,
	display_name = :display_name,
	department = :department,
	group_memberships = :group_memberships,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = :staff_id
	AND is_active = 1
LIMIT 1`

// UpdateAccount updates an active staff's account.
// A deactivated account must be re-activated
// before being updated.
//
// Input {userName: string, email: string, displayName: string, department: string, groupMembers: number}
func (env Env) UpdateAccount(p staff.Account) error {
	_, err := env.DB.NamedExec(stmtUpdateProfile, &p)
	if err != nil {
		logger.WithField("trace", "Env.UpdateAccount").Error(err)
		return err
	}

	return nil
}

const stmtDeactivate = `
UPDATE backyard.staff
  SET is_active = 0,
	deactivated_utc = UTC_TIMESTAMP()
WHERE staff_id = ?
  AND is_active = 1
LIMIT 1`

const stmtDeletePersonalToken = `
UPDATE oauth.access
	SET is_active = 0
WHERE created_by = ?`

// Deactivate a staff.
// Input {revokeVip: true | false}
func (env Env) Deactivate(id string) error {
	log := logger.WithField("trace", "Env.Deactivate")

	tx, err := env.DB.Beginx()
	if err != nil {
		log.Error(err)
		return err
	}

	// 1. Find the staff to deactivate.
	var account staff.Account
	if err := tx.Get(&account, stmtAccountByID, id); err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	if !account.IsActive {
		_ = tx.Rollback()
		return nil
	}

	// 2. Deactivate the staff
	_, err = tx.Exec(stmtDeactivate, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()

		return err
	}

	// 3. Remove personal tokens
	_, err = tx.Exec(stmtDeletePersonalToken, account.UserName)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

const stmtActivate = `
UPDATE backyard.staff
  SET is_active = 1,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = ?
  AND is_active = 0
LIMIT 1`

// Activate reinstate an deactivated account.
func (env Env) Activate(id string) error {
	_, err := env.DB.Exec(stmtActivate, id)

	if err != nil {
		logger.WithField("trace", "ActivateStaff").Error(err)

		return err
	}

	return nil
}
