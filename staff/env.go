package staff

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
)

var staffLogger = log.WithField("package", "staff")

// Env interact with user data
type Env struct {
	DB *sql.DB
}

// Auth perform authentication by user name and password
// POST /staff/auth
func (env Env) Auth(l Login) (Account, error) {
	query := `
	SELECT id AS id,
		username AS userName,
		email,
		display_name AS displayName,
		department AS department,
		group_memberships AS groups,
		vip_uuid AS myftId
	FROM backyard.staff
	WHERE (username, password) = (?, UNHEX(MD5(?)))
		AND is_active = 1
	LIMIT 1`

	var a Account
	err := env.DB.QueryRow(query, l.UserName, l.Password).Scan(
		&a.ID,
		&a.Email,
		&a.UserName,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
		&a.MyftID,
	)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "Auth",
			"table": "backyard.staff",
		}).Error(err)

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

// Profile retrieves all of a user's data.
// This is used by both an administrator or the user itself
// GET /user/profile
// GET /staff/profile
func (env Env) Profile(userName string) (Profile, error) {
	query := `
	SELECT staff.id AS id,
		staff.username AS userName,
		staff.email AS email,
		staff.is_active AS isActive,
		staff.display_name AS displayName,
		staff.department AS department,
		staff.group_memberships AS groupMembers,
		staff.vip_uuid AS vipId,
		myft.email AS vipEmail,
		staff.created_utc AS createdAt,
		staff.deactivated_utc AS deactivatedAt,
		staff.updated_utc AS updatedAt,
		staff.last_login_utc AS lastLoginAt,
		INET6_NTOA(staff.last_login_ip) AS lastLoginIp
  	FROM backyard.staff AS staff
    	LEFT JOIN cmstmp01.userinfo AS myft
		ON staff.vip_uuid = myft.user_id
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
		&p.MyftID,
		&p.MyftEmail,
		&p.CreatedAt,
		&p.DeactiviateAt,
		&p.UpdatedAt,
		&p.LastLoginAt,
		&p.LastLoginIP,
	)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "Profile",
			"table": "backyard.staff",
		}).Error(err)

		return p, err
	}

	return p, nil
}

// UpdateName allows a user to change its display name.
// PATCH /user/display-name
func (env Env) UpdateName(p Profile) error {
	query := `
	UPDATE backyard.staff
		SET display_name = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, p.DisplayName, p.UserName)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdateName",
			"table": "backyard.staff",
		}).Error(err)
		return err
	}

	return nil
}

// UpdateEmail allows a user to udpate its email address.
// PATH /user/email
func (env Env) UpdateEmail(p Profile) error {
	query := `
	UPDATE backyard.staff
		SET email = ?
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, p.DisplayName, p.UserName)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdateEmail",
			"table": "backyard.staff",
		}).Error(err)
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
func (env Env) UpdatePassword(p Password) error {
	// Verify user's old password
	matched, err := env.verifyPassword(p.UserName, p.Old)

	if err != nil {
		return err
	}

	if !matched {
		return errors.New("wrong password")
	}

	err = env.changePassword(p.UserName, p.New)

	if err != nil {
		return err
	}

	return nil
}
