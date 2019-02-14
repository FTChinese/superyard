package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/ftchinese/backyard-api/staff"
)

// StaffEnv interact with user data
type StaffEnv struct {
	DB *sql.DB
}

func (env StaffEnv) exists(col, value string) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS(
		SELECT *
		FROM backyard.staff
		WHERE %s = ?
	) AS alreadyExists`, col)

	var exists bool

	err := env.DB.QueryRow(query, value).Scan(&exists)

	if err != nil {
		logger.WithField("trace", "exists").Error(err)

		return false, err
	}

	return exists, nil
}

// NameExists checks if name exists in the user_name column of backyard.staff table.
func (env StaffEnv) NameExists(name string) (bool, error) {
	return env.exists(
		tableStaff.colName(),
		name)
}

// EmailExists checks if an email address exists in the email column of backyard.staff table.
func (env StaffEnv) EmailExists(email string) (bool, error) {
	return env.exists(
		tableStaff.colEmail(),
		email)
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
	LIMIT 1`, stmtStaffAccount, col, activeStmt)

	var a staff.Account
	err := env.DB.QueryRow(query, value).Scan(
		&a.ID,
		&a.Email,
		&a.UserName,
		&a.IsActive,
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

// LoadAccountByName retrieves a staff's account data by name.
func (env StaffEnv) LoadAccountByName(name string, active bool) (staff.Account, error) {
	return env.loadAccount(
		tableStaff.colName(),
		name,
		active)
}

// LoadAccountByEmail retrieves a staff's account data by email
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
		SET display_name = ?
	WHERE user_name = ?
		AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, displayName, userName)

	if err != nil {
		logger.WithField("trace", "UpdateName").Error(err)
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
	WHERE user_name = ?
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
func (env StaffEnv) Profile(userName string) (staff.Profile, error) {
	query := fmt.Sprintf(`
	%s
	WHERE user_name = ?
	LIMIT 1`, stmtStaffProfile)

	var p staff.Profile
	err := env.DB.QueryRow(query, userName).Scan(
		&p.ID,
		&p.Email,
		&p.UserName,
		&p.IsActive,
		&p.DisplayName,
		&p.Department,
		&p.GroupMembers,
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
