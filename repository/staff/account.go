package staff

import (
	"gitlab.com/ftchinese/backyard-api/models/employee"
)

func (env Env) Create(a employee.Account) error {
	_, err := env.DB.NamedExec(stmtInsertEmployee, &a)

	if err != nil {
		logger.WithField("trace", "Env.CreateAccount").Error(err)
		return err
	}

	return nil
}

// Deactivate a staff.
// Input {revokeVip: true | false}
func (env Env) Deactivate(id string, revokeVIP bool) error {
	log := logger.WithField("trace", "Env.Deactivate")

	tx, err := env.DB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	// 1. Deactivate a staff's account.
	_, err = tx.Exec(stmtDeactivate, id)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return err
	}

	// 2. Revoke VIP granted to all ftc accounts associated with this staff.
	if revokeVIP {
		_, err := tx.Exec(stmtRevokeVIP, id)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return err
		}
	}

	// 3. Delete myft accounts associated with this staff.
	_, err = tx.Exec(stmtDeleteMyft, id)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return err
	}

	// 4. Delete all access tokens to next-ap created by this user.
	_, err = tx.Exec(stmtDeletePersonalToken, id)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

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
