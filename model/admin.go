package model

import (
	"database/sql"
	"fmt"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/user"
	"gitlab.com/ftchinese/backyard-api/util"
)

type AdminEnv struct {
	DB *sql.DB
}

func (env AdminEnv) exists(col, value string) (bool, error) {
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
func (env AdminEnv) NameExists(name string) (bool, error) {
	return env.exists(
		tableStaff.colName(),
		name)
}

// EmailExists checks if an email address exists in the email column of backyard.staff table.
func (env AdminEnv) EmailExists(email string) (bool, error) {
	return env.exists(
		tableStaff.colEmail(),
		email)
}

// Create a new staff and generate a random password.
// The password is returned so that you could send it to user's email.
func (env AdminEnv) CreateAccount(a staff.Account, password string) error {

	query := `
	INSERT INTO backyard.staff
      SET user_name = ?,
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
		logger.WithField("trace", "CreateAccount").Error(err)
		return err
	}

	return nil
}

// ListAccounts list all staff with pagination support.
func (env AdminEnv) ListAccounts(p util.Pagination) ([]staff.Account, error) {
	query := fmt.Sprintf(`
	%s
	ORDER BY id ASC
	LIMIT ? OFFSET ?`, stmtStaffAccount)

	rows, err := env.DB.Query(
		query,
		p.RowCount,
		p.Offset())

	if err != nil {
		logger.
			WithField("location", "Query staff roster").
			Error(err)

		return nil, err
	}
	defer rows.Close()

	accounts := make([]staff.Account, 0)
	for rows.Next() {
		var a staff.Account

		err := rows.Scan(
			&a.ID,
			&a.UserName,
			&a.Email,
			&a.DisplayName,
			&a.Department,
			&a.GroupMembers,
		)

		if err != nil {
			logger.
				WithField("trace", "ListAccounts").
				Error(err)

			continue
		}

		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		logger.
			WithField("trace", "ListAccounts").
			Error(err)

		return accounts, err
	}

	return accounts, nil
}

// UpdateAccount updates a staff's profile by administrator
func (env AdminEnv) UpdateAccount(userName string, a staff.Account) error {
	query := `
	UPDATE backyard.staff
	SET user_name = ?,
		email = ?,
		display_name = ?,
		department = ?,
		group_memberships = ?,
		updated_utc = UTC_TIMESTAMP()
	WHERE user_name = ?
		AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query,
		a.UserName,
		a.Email,
		a.DisplayName,
		a.Department,
		a.GroupMembers,
		userName,
	)

	if err != nil {
		logger.WithField("trace", "UpdateAccount").Error(err)
		return err
	}

	return nil
}

// RemoveStaff deactivates a staff's account.
func (env AdminEnv) RemoveStaff(userName string, revokeVIP bool) error {
	tx, err := env.DB.Begin()

	// 1. Deactivate a staff's account.
	query := `
    UPDATE backyard.staff
	  SET is_active = 0,
	  	deactivated_utc = UTC_TIMESTAMP()
    WHERE user_name = ?
      AND is_active = 1
	LIMIT 1`

	_, err = tx.Exec(query, userName)

	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "RemoveStaff").Error(err)
	}

	// 2. Revoke VIP granted to all ftc accounts associated with this staff.
	if revokeVIP {
		query = `
	UPDATE backyard.staff_myft AS s
		LEFT JOIN cmstmp01.userinfo AS u
		ON s.myft_id = u.user_id
	SET is_vip = 0
	WHERE s.staff_name = ?`

		_, err = tx.Exec(query, userName)
		if err != nil {
			_ = tx.Rollback()
			logger.WithField("trace", "RemoveStaff").Error(err)
		}
	}

	// 3. Delete myft accounts associated with this staff.
	query = `
	DELETE FROM backyard.staff_myft
    WHERE staff_name = ?`

	_, err = tx.Exec(query, userName)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "RemoveStaff").Error(err)
	}

	// 4. Delete all access tokens to next-ap created by this user.
	query = `
	UPDATE oauth.access
		SET is_active = 0
	WHERE created_by = ?`

	_, err = tx.Exec(query, userName)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "RemoveStaff").Error(err)
	}

	if err := tx.Commit(); err != nil {
		logger.WithField("trace", "RemoveStaff").Error(err)
		return err
	}

	return nil
}

// ActivateStaff reuses a previously removed staff account
func (env AdminEnv) ActivateStaff(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 1
    WHERE user_name = ?
      AND is_active = 0
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		logger.WithField("trace", "ActivateStaff").Error(err)
		return err
	}

	return nil
}

// ListVIP list all vip account on ftchinese.com
func (env AdminEnv) ListVIP() ([]user.User, error) {

	query := fmt.Sprintf(`
	%s
	WHERE is_vip = 1`, stmtUser)

	rows, err := env.DB.Query(query)

	if err != nil {
		logger.WithField("trace", "ListVIP").Error(err)

		return nil, err
	}
	defer rows.Close()

	vips := make([]user.User, 0)
	for rows.Next() {
		var u user.User

		err := rows.Scan(
			&u.UserID,
			&u.UnionID,
			&u.Email,
			&u.UserName,
			&u.IsVIP,
		)

		if err != nil {
			logger.WithField("trace", "ListVIP").Error(err)

			continue
		}

		vips = append(vips, u)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListVIP").Error(err)

		return vips, err
	}

	return vips, nil
}

func (env AdminEnv) updateVIP(myftID string, isVIP bool) error {
	query := `
	UPDATE cmstmp01.userinfo
      SET is_vip = ?
    WHERE user_id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, isVIP, myftID)

	if err != nil {
		logger.WithField("location", "Grant vip to a ftc account")

		return err
	}

	return nil
}

// GrantVIP set a ftc account as vip
func (env AdminEnv) GrantVIP(myftID string) error {
	return env.updateVIP(myftID, true)
}

// RevokeVIP removes vip status from a ftc account
func (env AdminEnv) RevokeVIP(myftID string) error {
	return env.updateVIP(myftID, false)
}