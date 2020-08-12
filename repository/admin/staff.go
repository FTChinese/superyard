package admin

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/staff"
)

// CreateStaff creates a new staff account
func (env Env) CreateStaff(su staff.SignUp) error {
	_, err := env.DB.NamedExec(
		staff.StmtCreateAccount,
		su)

	if err != nil {
		logger.WithField("trace", "Env.CreateAccount").Error(err)
		return err
	}

	return nil
}

// AccountByID retrieves staff account by
// email column.
func (env Env) AccountByID(id string) (staff.Account, error) {
	var a staff.Account

	if err := env.DB.Get(&a, staff.StmtAccountByID, id); err != nil {
		return staff.Account{}, err
	}

	return a, nil
}

// AccountByName loads an account when by name
// is submitted to request a password reset letter.
func (env Env) AccountByName(name string) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, staff.StmtAccountByName, name)

	if err != nil {
		logger.WithField("trace", "Env.AccountByName").Error(err)

		return staff.Account{}, err
	}

	return a, err
}

func (env Env) ListStaff(p gorest.Pagination) ([]staff.Account, error) {
	accounts := make([]staff.Account, 0)

	err := env.DB.Select(&accounts,
		staff.ListAccounts,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListStaff").Error(err)

		return accounts, err
	}

	return accounts, nil
}

// UpdateAccount updates an active staff's account.
// A deactivated account must be re-activated
// before being updated.
//
// Input
// {
//   userName: string,
//   email: string,
//   displayName: string,
//   department: string,
//   groupMembers: number
//  }
func (env Env) UpdateAccount(p staff.Account) error {
	_, err := env.DB.NamedExec(staff.StmtUpdateAccount, &p)
	if err != nil {
		logger.WithField("trace", "Env.UpdateAccount").Error(err)
		return err
	}

	return nil
}

// StaffProfile loads a staff's profile.
func (env Env) StaffProfile(id string) (staff.Profile, error) {
	var p staff.Profile

	err := env.DB.Get(&p, staff.StmtProfile, id)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveProfile").Error(err)

		return p, err
	}

	return p, nil
}

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
	if err := tx.Get(&account, staff.StmtAccountByID, id); err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	if !account.IsActive {
		_ = tx.Rollback()
		return nil
	}

	// 2. Deactivate the staff
	_, err = tx.Exec(staff.StmtDeactivate, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()

		return err
	}

	// 3. Remove personal tokens
	_, err = tx.Exec(
		staff.StmtDeletePersonalKey,
		account.UserName,
	)
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

// Activate reinstate an deactivated account.
func (env Env) Activate(id string) error {
	_, err := env.DB.Exec(staff.StmtActivate, id)

	if err != nil {
		logger.WithField("trace", "ActivateStaff").Error(err)

		return err
	}

	return nil
}
