package model

import (
	"database/sql"
	"fmt"
	"gitlab.com/ftchinese/backyard-api/staff"
)

// Env interact with user data
type StaffEnv struct {
	DB *sql.DB
}

// UpdateLoginHistory saves user login footprint after successfully authenticated.
func (env StaffEnv) UpdateLoginHistory(l staff.Login, ip string) error {
	query := `
    UPDATE backyard.staff
      SET last_login_utc = UTC_TIMESTAMP(),
        last_login_ip = IFNULL(INET6_ATON(?), last_login_ip)
    WHERE username = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, ip, l.UserName)

	if err != nil {
		logger.WithField("trace", "UpdateLoginHistory").Error(err)
		return err
	}

	return nil
}

// LoadAccount gets an account by user name.
// Use `activeOnly` to limit active staff only or all.
func (env StaffEnv) loadAccount(col, value string, activeOnly bool) (staff.Account, error) {
	var activeStmt string
	if activeOnly {
		activeStmt = "AND is_active = 1"
	}
	query := fmt.Sprintf(`
	%s
	WHERE %s = ?
		%s	
	LIMIT 1`, stmtStaffAccount, string(col), activeStmt)

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

func (env StaffEnv) LoadAccountByName(name string, active bool) (staff.Account, error) {
	return env.loadAccount(
		tableStaff.colName(),
		name,
		active)
}

func (env StaffEnv) LoadAccountByEmail(email string, active bool) (staff.Account, error) {
	return env.loadAccount(
		tableStaff.colEmail(),
		email,
		active)
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
	query := fmt.Sprintf(`
	%s
	WHERE username = ?
	LIMIT 1`, stmtStaffProfile)

	var p staff.Profile
	err := env.DB.QueryRow(query, userName).Scan(
		&p.ID,
		&p.Email,
		&p.UserName,
		&p.DisplayName,
		&p.Department,
		&p.GroupMembers,
		&p.IsActive,
		&p.CreatedAt,
		&p.DeactivatedAt,
		&p.UpdatedAt,
		&p.LastLoginAt,
		&p.LastLoginIP,
	)

	if err != nil {
		logger.WithField("trace", "Profile").Error(err)
		return p, err
	}

	return p, nil
}