package admin

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/staff"
	"log"
)

// CreateStaff creates a new staff account
func (env Env) CreateStaff(su staff.SignUp) error {
	_, err := env.db.NamedExec(
		staff.StmtCreateAccount,
		su)

	if err != nil {
		return err
	}

	return nil
}

// AccountByID retrieves staff account by
// email column.
func (env Env) AccountByID(id string) (staff.Account, error) {
	var a staff.Account

	if err := env.db.Get(&a, staff.StmtAccountByID, id); err != nil {
		return staff.Account{}, err
	}

	return a, nil
}

// AccountByName loads an account when by name
// is submitted to request a password reset letter.
func (env Env) AccountByName(name string) (staff.Account, error) {
	var a staff.Account
	err := env.db.Get(&a, staff.StmtAccountByName, name)

	if err != nil {
		return staff.Account{}, err
	}

	return a, err
}

func (env Env) countStaff() (int64, error) {
	var count int64

	err := env.db.Get(&count, staff.StmtCountStaff)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listStaff(p gorest.Pagination) ([]staff.Account, error) {
	accounts := make([]staff.Account, 0)

	err := env.db.Select(&accounts,
		staff.StmtListAccounts,
		p.Limit,
		p.Offset())

	if err != nil {
		return accounts, err
	}

	return accounts, nil
}

func (env Env) ListStaff(p gorest.Pagination) (staff.AccountList, error) {
	countCh := make(chan int64)
	listCh := make(chan staff.AccountList)

	go func() {
		defer close(countCh)

		n, err := env.countStaff()
		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)

		list, err := env.listStaff(p)
		if err != nil {
			log.Print(err)
		}

		listCh <- staff.AccountList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return staff.AccountList{}, listResult.Err
	}

	return staff.AccountList{
		Total:      count,
		Pagination: p,
		Data:       listResult.Data,
		Err:        nil,
	}, nil
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
	_, err := env.db.NamedExec(staff.StmtUpdateAccount, &p)
	if err != nil {
		return err
	}

	return nil
}

// StaffProfile loads a staff's profile.
func (env Env) StaffProfile(id string) (staff.Profile, error) {
	var p staff.Profile

	err := env.db.Get(&p, staff.StmtProfile, id)

	if err != nil {
		return p, err
	}

	return p, nil
}

// Deactivate a staff.
// Input {revokeVip: true | false}
func (env Env) Deactivate(id string) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	// 1. Find the staff to deactivate.
	var account staff.Account
	if err := tx.Get(&account, staff.StmtAccountByID, id); err != nil {
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
		_ = tx.Rollback()

		return err
	}

	// 3. Remove personal tokens
	_, err = tx.Exec(
		staff.StmtDeletePersonalKey,
		account.UserName,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Activate reinstate an deactivated account.
func (env Env) Activate(id string) error {
	_, err := env.db.Exec(staff.StmtActivate, id)

	if err != nil {
		return err
	}

	return nil
}
