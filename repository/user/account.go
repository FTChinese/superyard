package user

import (
	"github.com/FTChinese/superyard/pkg/staff"
)

// AccountByID retrieves staff account by id column.
func (env Env) AccountByID(id string) (staff.Account, error) {
	var a staff.Account

	if err := env.DB.Get(&a, staff.StmtActiveAccountByID, id); err != nil {
		return staff.Account{}, err
	}

	return a, nil
}

// AccountByEmail loads an account when a email
// is submitted to request a password reset letter.
func (env Env) AccountByEmail(email string) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, staff.StmtActiveAccountByEmail, email)

	if err != nil {
		return staff.Account{}, err
	}

	return a, err
}

func (env Env) AddID(a staff.Account) error {

	_, err := env.DB.NamedExec(staff.StmtAddID, a)

	if err != nil {
		return err
	}

	return nil
}

// SetEmail sets the email column is missing.
func (env Env) SetEmail(a staff.Account) error {
	_, err := env.DB.NamedExec(staff.StmtSetEmail, a)

	if err != nil {
		return err
	}

	return nil
}

// UpdateDisplayName changes display name.
func (env Env) UpdateDisplayName(a staff.Account) error {
	_, err := env.DB.NamedExec(staff.StmtUpdateDisplayName, a)

	if err != nil {
		return err
	}

	return nil
}

// RetrieveProfile loads a staff's profile.
func (env Env) RetrieveProfile(id string) (staff.Profile, error) {
	var p staff.Profile

	err := env.DB.Get(&p, staff.StmtActiveProfile, id)

	if err != nil {
		return p, err
	}

	return p, nil
}
