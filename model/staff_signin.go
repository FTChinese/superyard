package model

import (
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
)

// Auth perform authentication by user name and password
// POST /staff/auth
func (env StaffEnv) Auth(l staff.Login) (staff.Account, error) {
	// Verify password
	matched, err := env.isPasswordMatched(l.UserName, l.Password)

	// User might not be found
	if err != nil {
		return staff.Account{}, err
	}

	// Password is incorrect
	if !matched {
		return staff.Account{}, util.ErrWrongPassword
	}

	a, err := env.FindAccountByName(l.UserName, true)

	if err != nil {
		return a, err
	}

	go env.updateLoginHistory(l)

	return a, nil
}

// UpdateLoginHistory saves user login footprint after successfully authenticated.
func (env StaffEnv) updateLoginHistory(l staff.Login) error {
	query := `
    UPDATE backyard.staff
      SET last_login_utc = UTC_TIMESTAMP(),
        last_login_ip = IFNULL(INET6_ATON(?), last_login_ip)
    WHERE username = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, l.UserIP, l.UserName)

	if err != nil {
		logger.WithField("location", "Update login history").Error(err)
		return err
	}

	return nil
}
