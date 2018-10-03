package staff

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var staffLogger = log.WithField("package", "staff")

type sqlCol string

const (
	colUserName sqlCol = "username"
	colEmail    sqlCol = "email"
	stmtAccount string = `
	SELECT id AS id,
		username AS userName,
		IFNULL(email, '') AS email,
		IFNULL(display_name, '') AS displayName,
		IFNULL(department, '') AS department,
		group_memberships AS groups
	FROM backyard.staff`
)

// Env interact with user data
type Env struct {
	DB *sql.DB
}

func (env Env) exists(col sqlCol, value string) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS(
		SELECT *
		FROM backyard.staff
		WHERE %s = ?
	) AS alreadyExists`, string(col))

	var exists bool

	err := env.DB.QueryRow(query, value).Scan(&exists)

	if err != nil {
		staffLogger.
			WithField("location", "staff exists").
			Error(err)

		return false, err
	}

	return exists, nil
}

// StaffNameExists checks if name exists in the username column of backyard.staff table.
func (env Env) StaffNameExists(name string) (bool, error) {
	return env.exists(colUserName, name)
}

// StaffEmailExists checks if an email address exists in the email column of backyard.staff table.
func (env Env) StaffEmailExists(email string) (bool, error) {
	return env.exists(colEmail, email)
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
		staffLogger.WithField("location", "Staff authentication").Error(err)

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
		staffLogger.WithFields(log.Fields{
			"func":  "UpdateLoginHistory",
			"table": "backyard.staff",
		}).Error(err)

		return err
	}

	return nil
}

func (env Env) findAccount(col sqlCol, value string) (Account, error) {
	query := fmt.Sprintf(`
	%s
	WHERE %s = ?
		AND is_active = 1
	LIMIT 1`, stmtAccount, string(col))

	var a Account
	err := env.DB.QueryRow(query, value).Scan(
		&a.ID,
		&a.Email,
		&a.UserName,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
	)

	if err != nil {
		staffLogger.WithField("location", "Find account by username or email").Error(err)

		return a, err
	}

	return a, nil
}

// Profile retrieves all of a user's data.
// This is used by both an administrator or the user itself
// GET /user/profile
// GET /staff/profile
func (env Env) Profile(userName string) (Profile, error) {
	query := `
	SELECT id AS id,
		username AS userName,
		email,
		is_active AS isActive,
		display_name AS displayName,
		department AS department,
		group_memberships AS groupMembers,
		created_utc AS createdAt,
		deactivated_utc AS deactivatedAt,
		updated_utc AS updatedAt,
		last_login_utc AS lastLoginAt,
		INET6_NTOA(staff.last_login_ip) AS lastLoginIp
  	FROM backyard.staff
	WHERE username = ?
	LIMIT 1`

	var p Profile
	err := env.DB.QueryRow(query, userName).Scan(
		&p.ID,
		&p.UserName,
		&p.Email,
		&p.IsActive,
		&p.DisplayName,
		&p.Department,
		&p.GroupMembers,
		&p.CreatedAt,
		&p.DeactiviateAt,
		&p.UpdatedAt,
		&p.LastLoginAt,
		&p.LastLoginIP,
	)

	if err != nil {
		staffLogger.WithField("location", "Retrieving staff profile").Error(err)

		return p, err
	}

	return p, nil
}

// UpdateName allows a user to change its display name.
// PATCH /user/display-name
func (env Env) UpdateName(userName string, displayName string) error {
	query := `
	UPDATE backyard.staff
		SET display_name = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, displayName, userName)

	if err != nil {
		staffLogger.WithField("location", "Updating staff name").Error(err)
		return err
	}

	return nil
}

// UpdateEmail allows a user to udpate its email address.
// PATH /user/email
func (env Env) UpdateEmail(userName string, email string) error {
	query := `
	UPDATE backyard.staff
		SET email = ?
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, email, userName)

	if err != nil {
		staffLogger.WithField("location", "").Error(err)
		return err
	}

	return nil
}

// UpdatePassword allows a user to change password via email or username
// POST /user/password
func (env Env) changePassword(userName string, password string) error {
	query := `
	UPDATE backyard.staff
		SET password = UNHEX(MD5(?)),
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
		AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, password, userName)

	if err != nil {
		staffLogger.WithField("location", "Update backyard.staff password").Error(err)
		return err
	}

	legacyQuery := `
	UPDATE cmstmp01.managers
		SET password = MD5(?)
	WHERE %s = ?
	LIMIT 1`

	_, err = env.DB.Exec(legacyQuery, password, userName)

	if err != nil {
		staffLogger.WithField("location", "Update cmstmp01.managers password").Error(err)
		return err
	}

	return nil
}

// verifyPassword is used when a logged in user tries to change its password
func (env Env) verifyPassword(userName string, password string) (bool, error) {
	query := `
	SELECT password = UNHEX(MD5(?)) AS matched
	FROM backyard.staff
	WHERE username = ?
	LIMIT 1`

	var matched bool
	err := env.DB.QueryRow(query, password, userName).Scan(&matched)

	if err != nil {
		staffLogger.WithField("location", "Verify password").Error(err)

		return matched, err
	}

	return matched, nil
}

// UpdatePassword allows user to change password in its settings.
func (env Env) UpdatePassword(userName string, p Password) error {
	// Verify user's old password
	matched, err := env.verifyPassword(userName, p.Old)

	if err != nil {
		return err
	}

	if !matched {
		return errors.New("wrong password")
	}

	err = env.changePassword(userName, p.New)

	if err != nil {
		return err
	}

	return nil
}
