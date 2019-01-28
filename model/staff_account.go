package model

import (
	"fmt"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
)

func (env StaffEnv) exists(col sqlCol, value string) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS(
		SELECT *
		FROM backyard.staff
		WHERE %s = ?
	) AS alreadyExists`, string(col))

	var exists bool

	err := env.DB.QueryRow(query, value).Scan(&exists)

	if err != nil {
		logger.
			WithField("location", "staff exists").
			Error(err)

		return false, err
	}

	return exists, nil
}

// StaffNameExists checks if name exists in the username column of backyard.staff table.
func (env StaffEnv) StaffNameExists(name string) (bool, error) {
	return env.exists(colUserName, name)
}

// StaffEmailExists checks if an email address exists in the email column of backyard.staff table.
func (env StaffEnv) StaffEmailExists(email string) (bool, error) {
	return env.exists(colEmail, email)
}

// Create a new staff and generate a random password.
// The password is returned so that you could send it to user's email.
func (env StaffEnv) CreateAccount(a staff.Account, password string) error {

	query := `
	INSERT INTO backyard.staff
      SET username = ?,
        email = ?,
        password = UNHEX(MD5(?)),
        display_name = ?,
        department = ?,
		group_memberships = ?`

	_, err := env.DB.Exec(query,
		a.UserName,
		a.Email,
		password,
		a.DisplayName,
		a.Department,
		a.GroupMembers,
	)

	if err != nil {
		logger.
			WithField("location", "Inserting new staff").
			Error(err)

		return err
	}

	return nil
}

// FindAccount gets an account by user name.
// Use `activeOnly` to limit active staff only or all.
func (env StaffEnv) findAccount(col sqlCol, value string, activeOnly bool) (staff.Account, error) {
	var activeStmt string
	if activeOnly {
		activeStmt = "AND is_active = 1"
	}
	query := fmt.Sprintf(`
	SELECT id AS id,
		username AS userName,
		IFNULL(email, '') AS email,
		display_name AS displayName,
		department AS department,
		group_memberships AS groups
	FROM backyard.staff
	WHERE %s = ?
		%s	
	LIMIT 1`, string(col), activeStmt)

	var a staff.Account
	err := env.DB.QueryRow(query, value).Scan(
		&a.ID,
		&a.UserName,
		&a.Email,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
	)

	if err != nil {
		logger.WithField("location", "Staff authentication").Error(err)

		return a, err
	}

	return a, nil
}

func (env StaffEnv) FindAccountByName(name string, active bool) (staff.Account, error) {
	return env.findAccount(colUserName, name, active)
}

func (env StaffEnv) FindAccountByEmail(email string, active bool) (staff.Account, error) {
	return env.findAccount(colEmail, email, active)
}

// UpdateName allows a user to change its display name.
// PATCH /user/display-name
func (env StaffEnv) UpdateName(userName string, displayName string) error {
	query := `
	UPDATE backyard.staff
		SET display_name = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
		AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, displayName, userName)

	if err != nil {
		logger.WithField("location", "Updating staff name").Error(err)
		return err
	}

	return nil
}

// UpdateEmail allows a user to udpate its email address.
// PATH /user/email
func (env StaffEnv) UpdateEmail(userName string, email string) error {
	query := `
	UPDATE backyard.staff
		SET email = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
		AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, email, userName)

	if err != nil {
		logger.WithField("location", "").Error(err)
		return err
	}

	return nil
}

// Profile retrieves all of a user's data.
// This is used by both an administrator or the user itself
// GET /user/profile
// GET /staff/profile
func (env StaffEnv) Profile(userName string) (staff.Profile, error) {
	query := `
	SELECT id AS id,
		username AS userName,
		IFNULL(email, '') AS email,
		is_active AS isActive,
		IFNULL(display_name, '') AS displayName,
		IFNULL(department, '') AS department,
		group_memberships AS groupMembers,
		created_utc AS createdAt,
		IFNULL(deactivated_utc, '') AS deactivatedAt,
		IFNULL(updated_utc, '') AS updatedAt,
		IFNULL(last_login_utc, '') AS lastLoginAt,
		IFNULL(INET6_NTOA(staff.last_login_ip), '') AS lastLoginIp
  	FROM backyard.staff
	WHERE username = ?
	LIMIT 1`

	var p staff.Profile
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
		logger.WithField("location", "Retrieving staff profile").Error(err)

		return p, err
	}

	p.CreatedAt = util.ISO8601UTC.FromDatetime(p.CreatedAt, nil)
	if p.DeactiviateAt != "" {
		p.DeactiviateAt = util.ISO8601UTC.FromDatetime(p.DeactiviateAt, nil)
	}

	if p.UpdatedAt != "" {
		p.UpdatedAt = util.ISO8601UTC.FromDatetime(p.UpdatedAt, nil)
	}

	if p.LastLoginAt != "" {
		p.LastLoginAt = util.ISO8601UTC.FromDatetime(p.LastLoginAt, nil)
	}

	return p, nil
}