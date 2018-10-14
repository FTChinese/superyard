package staff

import (
	"strings"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Login specifies the the fields used for authentication
type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	UserIP   string `json:"userIp"`
}

// Sanitize removes leading and trailing space of each field
func (l *Login) Sanitize() {
	l.UserName = strings.TrimSpace(l.UserName)
	l.Password = strings.TrimSpace(l.Password)
	l.UserIP = strings.TrimSpace(l.UserIP)
}

// Auth perform authentication by user name and password
// POST /staff/auth
func (env Env) Auth(l Login) (Account, error) {
	// Verify password
	matched, err := env.isPasswordMatched(l.UserName, l.Password)

	// User might not be found
	if err != nil {
		return Account{}, err
	}

	// Password is incorrect
	if !matched {
		return Account{}, util.ErrWrongPassword
	}

	a, err := env.FindAccount(l.UserName, true)

	if err != nil {
		return a, err
	}

	go env.updateLoginHistory(l)

	return a, nil
}

// UpdateLoginHistory saves user login footprint after successfully authenticated.
func (env Env) updateLoginHistory(l Login) error {
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
