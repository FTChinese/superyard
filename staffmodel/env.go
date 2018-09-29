package staffmodel

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var staffLogger = log.WithFields(log.Fields{
	"package":  "staffmodel",
	"resource": "Env",
})

// SQLCol dynamically determines which SQL column is used
type SQLCol int

// SQL columns
const (
	ColUserName SQLCol = 0
	ColEmail    SQLCol = 1
)

func (col SQLCol) String() string {
	cols := [...]string{
		"username",
		"email",
	}

	return cols[col]
}

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

	go env.UpdateLoginHistory(l)

	return a, nil
}

// RetrieveAccount get a user's account by email or by username column
func (env Env) RetrieveAccount(col SQLCol, value string) (Account, error) {
	query := fmt.Sprintf(`
	SELECT id AS id,
		username AS userName,
		email,
		display_name AS displayName,
		department AS department,
		group_memberships AS groupMembers,
		vip_uuid AS myftId
	FROM backyard.staff
	WHERE %s = ?
	LIMIT 1`, col.String())

	var a Account
	err := env.DB.QueryRow(query, value).Scan(
		&a.ID,
		&a.Email,
		&a.UserName,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
		&a.MyftID,
	)

	if err != nil {
		staffLogger.WithField("location", "Retrieve a staff account").Error(err)

		return a, err
	}

	return a, nil
}

// UpdateLoginHistory saves user login footprint after successfully authenticated.
func (env Env) UpdateLoginHistory(l Login) error {
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

// Exists checks if a username or email exists
func (env Env) Exists(col SQLCol, value string) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS(
		SELECT *
		FROM backyard.staff
		WHERE %s = ?
	) AS alreadyExists`, col.String())

	var exists bool

	err := env.DB.QueryRow(query, value).Scan(&exists)

	if err != nil {
		staffLogger.
			WithField("func", "Exists").
			Error(err)

		return false, err
	}

	return exists, nil
}

// Create adds a new staff's profile
// POST /staff/new
func (env Env) Create(s Staff) error {
	query := `
	INSERT INTO backyard.staff
      SET username = ?,
        email = ?,
        password = UNHEX(MD5(?)),
        display_name = NULLIF(?, ''),
        department = NULLIF(?, ''),
		group_memberships = ?`

	_, err := env.DB.Exec(query,
		s.UserName,
		s.Email,
		s.Password,
		s.DisplayName,
		s.Department, s.GroupMembers)

	if err != nil {
		staffLogger.
			WithField("func", "CreateStaff").
			Error(err)

		return err
	}

	return nil
}

// Activate turns a staff as active if it was removed
// PUT /staff/new
func (env Env) Activate(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 1
    WHERE username = ?
      AND is_active = 0
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		staffLogger.
			WithField("func", "ActivateStaff").
			Error(err)

		return err
	}

	return nil
}

// Roster list all staff
// GET /staff/roster
func (env Env) Roster(page int, rowCount int) ([]Account, error) {
	offset := (page - 1) * rowCount
	query := `
	SELECT id AS id,
		username AS userName,
		display_name AS displayName,
		department AS department,
		group_memberships AS groupMembers,
		myft_id AS myftId
	FROM backyard.staff
	WHERE is_active = 1
	ORDER BY id ASC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(query, rowCount, offset)

	var items []Account

	if err != nil {
		staffLogger.
			WithField("func", "StaffRoster").
			Error(err)

		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Account

		err := rows.Scan(
			&item.ID,
			&item.UserName,
			&item.DisplayName,
			&item.Department,
			&item.GroupMembers,
			&item.MyftID,
		)

		if err != nil {
			staffLogger.
				WithField("func", "StaffRoster").
				Error(err)

			continue
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		staffLogger.
			WithField("func", "StaffRoster").
			Error(err)

		return items, err
	}

	return items, nil
}

// Deactivate removes a staff
// DELETE /staff/profile
func (env Env) Deactivate(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 0
    WHERE userName = ?
      AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		staffLogger.
			WithField("func", "DeactivateStaff").
			Error(err)

		return err
	}

	return nil
}

// RemoveVIP set vip to false for all ftc accounts associated with a staff
// This should be perfomed when you deactivate a staff's account.
func (env Env) RemoveVIP(userName string) error {
	query := `
	UPDATE backyard.staff_myft AS s
		LEFT JOIN cmstmp01.userinfo AS u
		ON s.myft_id = u.user_id
	SET isvip = 0
	WHERE s.staff_name = ?
		AND u.isvip = 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		staffLogger.WithField("location", "remove vip status of a staff").Error(err)

		return nil
	}

	return nil
}

// UpdateProfile updates a user's profile by administrator
// PATCH /staff/profile
func (env Env) UpdateProfile(p Profile) error {
	query := `
	UPDATE backyard.staff
	SET username = ?,
		email = ?,
		display_name = NULLIF(?, ''),
		department = NULLIF(?, ''),
		group_memberships = ?,
		updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
		AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query,
		p.UserName,
		p.Email,
		p.DisplayName,
		p.Department,
		p.GroupMembers,
		p.UserName,
	)

	if err != nil {
		staffLogger.
			WithField("func", "UpdateProfile").
			Error(err)

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
func (env Env) UpdatePassword(col SQLCol, colValue string, password string) error {
	query := fmt.Sprintf(`
	UPDATE backyard.staff
		SET password = UNHEX(MD5(?)),
			updated_utc = UTC_TIMESTAMP()
	WHERE %s = ?
		AND is_active = 1
	LIMIT 1`, col.String())

	_, err := env.DB.Exec(query, password, colValue)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdatePassword",
			"table": "backyard.staff",
		}).Error(err)
		return err
	}

	legacyQuery := fmt.Sprintf(`
	UPDATE cmstmp01.managers
		SET password = MD5(?)
	WHERE %s = ?
	LIMIT 1`, col.String())

	_, err = env.DB.Exec(legacyQuery, password, colValue)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdatePassword",
			"table": "cmstmp01.managers",
		}).Error(err)
		return err
	}

	return nil
}
