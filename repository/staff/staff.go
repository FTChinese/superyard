package staff

import (
	"github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/models/util"
)

func (env Env) Create(a employee.Account) error {
	_, err := env.DB.NamedExec(stmtInsertEmployee, &a)

	if err != nil {
		logger.WithField("trace", "Env.CreateAccount").Error(err)
		return err
	}

	return nil
}

func (env Env) Load(col Column, value string) (employee.Profile, error) {
	var p employee.Profile

	err := env.DB.Get(&p, queryProfile(col), value)

	if err != nil {
		logger.WithField("trace", "Env.LoadProfile").Error(err)

		return p, err
	}

	return p, nil
}

func (env Env) List(p gorest.Pagination) ([]employee.Profile, error) {
	profiles := make([]employee.Profile, 0)

	err := env.DB.Select(&profiles,
		stmtListStaff,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.List").Error(err)

		return profiles, err
	}

	return profiles, nil
}

// UpdateAccount updates a staff's account.
//
//	PATCH /admin/accounts/{name}
//
// Input {userName: string, email: string, displayName: string, department: string, groupMembers: number}
func (env Env) Update(p employee.Profile) error {
	_, err := env.DB.NamedExec(stmtUpdateProfile, &p)
	if err != nil {
		logger.WithField("trace", "Env.UpdateAccount").Error(err)
		return err
	}

	return nil
}

// UpdateDisplayName allows a user to change its display name.
// PATCH /user/display-name
func (env Env) UpdateDisplayName(displayName, staffID string) error {
	_, err := env.DB.Exec(
		stmtUpdateName,
		displayName,
		staffID)

	if err != nil {
		logger.WithField("trace", "Env.UpdateDisplayName").Error(err)

		return err
	}

	return nil
}

// UpdateEmail allows a user to update its email address.
func (env Env) UpdateEmail(email, staffID string) error {
	_, err := env.DB.Exec(stmtUpdateEmail, email, staffID)

	if err != nil {
		logger.WithField("trace", "Env.UpdateEmail").Error(err)

		return err
	}

	return nil
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
func (env Env) UpdatePassword(p employee.Password, userName string) error {
	matched, err := env.VerifyPassword(employee.Login{
		UserName: userName,
		Password: p.Old,
	})

	if err != nil {
		return err
	}

	if !matched {
		return util.ErrWrongPassword
	}

	if err := env.changePassword(p.New, userName); err != nil {
		return err
	}

	return nil
}
