package user

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

const stmtLogin = stmt.StaffAccount + `
FROM backyard.staff AS s
WHERE (s.user_name, .s.password) = (?, UNHEX(MD5(?)))
	AND s.is_active = 1`

// Login verifies user name and password combination.
func (env Env) Login(l employee.Login) (employee.Account, error) {
	var a employee.Account
	err := env.DB.Get(&a, stmtLogin, l.UserName, l.Password)

	if err != nil {
		logger.WithField("trace", "Env.Login").Error(err)

		return a, err
	}

	return a, nil
}

const stmtUpdateLastLogin = `
UPDATE backyard.staff
SET last_login_utc = UTC_TIMESTAMP(),
	last_login_ip = IFNULL(INET6_ATON(?), last_login_ip)
WHERE user_name = ?
LIMIT 1`

// UpdateLastLogin saves user login footprint after successfully authenticated.
func (env Env) UpdateLastLogin(l employee.Login, ip string) error {
	_, err := env.DB.Exec(stmtUpdateLastLogin, ip, l.UserName)

	if err != nil {
		logger.WithField("trace", "Env.UpdateLastLogin").Error(err)

		return err
	}

	return nil
}
