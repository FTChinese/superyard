package user

import (
	"gitlab.com/ftchinese/superyard/models/staff"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

const accountByID = stmt.StaffAccount + `
FROM backyard.staff AS s
WHERE s.staff_id = ?
	AND s.is_active = 1
LIMIT 1`

// AccountByID retrieves staff account by
// email column.
func (env Env) AccountByID(id string) (staff.Account, error) {
	var a staff.Account

	if err := env.DB.Get(&a, accountByID, id); err != nil {
		return staff.Account{}, err
	}

	return a, nil
}

const stmtAccountByEmail = stmt.StaffAccount + `
FROM backyard.staff AS s
WHERE s.email = ?
	AND s.is_active = 1
LIMIT 1`

// AccountByEmail loads an account when a email
// is submitted to request a password reset letter.
func (env Env) AccountByEmail(email string) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, stmtAccountByEmail, email)

	if err != nil {
		logger.WithField("trace", "Env.AccountByEmail").Error(err)

		return staff.Account{}, err
	}

	return a, err
}

const stmtAddID = `
UPDATE backyard.staff
SET staff_id = :staff_id
WHERE user_name = :user_name
LIMIT 1`

func (env Env) AddID(a staff.Account) error {

	_, err := env.DB.NamedExec(stmtAddID, a)

	if err != nil {
		logger.WithField("trace", "Env.AddID").Error(err)
		return err
	}

	return nil
}

const stmtSetEmail = `
UPDATE backyard.staff
SET email = :email,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = :staff_id`

// SetEmail sets the email column is missing.
func (env Env) SetEmail(a staff.Account) error {
	_, err := env.DB.NamedExec(stmtSetEmail, a)

	if err != nil {
		logger.WithField("trace", "Env.SetEmail").Error(err)
		return err
	}

	return nil
}

const stmtDisplayName = `
UPDATE backyard.staff
SET display_name = :display_name,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = :staff_id`

// UpdateDisplayName changes display name.
func (env Env) UpdateDisplayName(a staff.Account) error {
	_, err := env.DB.NamedExec(stmtDisplayName, a)

	if err != nil {
		logger.WithField("trace", "Env.UpdateDisplayName").Error(err)
		return err
	}

	return nil
}

const stmtSelectProfile = stmt.StaffProfile + `
WHERE s.staff_id = ?
	AND s.is_active = 1
LIMIT 1`

// RetrieveProfile loads a staff's profile.
func (env Env) RetrieveProfile(id string) (staff.Profile, error) {
	var p staff.Profile

	err := env.DB.Get(&p, stmtSelectProfile, id)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveProfile").Error(err)

		return p, err
	}

	return p, nil
}
