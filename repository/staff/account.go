package staff

import "gitlab.com/ftchinese/backyard-api/models/employee"

func (env Env) Create(a employee.Account) error {
	_, err := env.DB.NamedExec(stmtCreateAccount, &a)

	if err != nil {
		logger.WithField("trace", "Env.CreateAccount").Error(err)
		return err
	}

	return nil
}

// RetrieveAccount retrieves staff account by
// email column.
func (env Env) RetrieveAccount(col employee.Column, val string) (employee.Account, error) {
	var a employee.Account

	if err := env.DB.Get(&a, QueryAccount(col), val); err != nil {
		return employee.Account{}, err
	}

	return a, nil
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
	var account employee.Account
	if err := tx.Get(&account, QueryAccount(employee.ColumnStaffID), id); err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
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

func (env Env) Activate(id string) error {
	_, err := env.DB.Exec(stmtActivate, id)

	if err != nil {
		logger.WithField("trace", "ActivateStaff").Error(err)

		return err
	}

	return nil
}

func (env Env) AddID(a employee.Account) error {

	_, err := env.DB.NamedExec(stmtAddID, a)

	if err != nil {
		logger.WithField("trace", "Env.AddID").Error(err)
		return err
	}

	return nil
}
