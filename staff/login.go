package staff

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
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
	query := fmt.Sprintf(`
	%s
	WHERE (username, password) = (?, UNHEX(MD5(?)))
		AND is_active = 1
	LIMIT 1`, stmtAccount)

	var a Account
	err := env.DB.QueryRow(query, l.UserName, l.Password).Scan(
		&a.ID,
		&a.Email,
		&a.UserName,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
	)

	if err != nil {
		logger.WithField("location", "Staff authentication").Error(err)

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
        last_login_ip = IFNULL(INET6_ATON(:?), last_login_ip)
    WHERE username = :?
	LIMIT 1`

	_, err := env.DB.Exec(query, l.UserIP, l.UserName)

	if err != nil {
		logger.WithFields(log.Fields{
			"func":  "UpdateLoginHistory",
			"table": "backyard.staff",
		}).Error(err)

		return err
	}

	return nil
}
