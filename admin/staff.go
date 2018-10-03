package admin

import (
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
)

func (env Env) createStaff(a staff.Account) (string, error) {
	password, err := util.RandomHex(4)

	if err != nil {
		adminLogger.WithField("location", "Creating password for new staff").Error(err)

		return "", err
	}

	query := `
	INSERT INTO backyard.staff
      SET username = ?,
        email = ?,
        password = UNHEX(MD5(?)),
        display_name = NULLIF(?, ''),
        department = NULLIF(?, ''),
		group_memberships = ?`

	_, err = env.DB.Exec(query,
		a.UserName,
		a.Email,
		password,
		a.DisplayName,
		a.Department,
		a.GroupMembers,
	)

	if err != nil {
		adminLogger.
			WithField("location", "Inserting new staff").
			Error(err)

		return "", err
	}

	return password, nil
}

// NewStaff creates a new account for a staff
// After the account is created, you should send the password to this its email address.
func (env Env) NewStaff(a staff.Account) error {
	pass, err := env.createStaff(a)

	if err != nil {
		return err
	}

	err = a.SendPassword(pass, newStaffLetterURL)

	if err != nil {
		return err
	}

	return nil
}

// StaffRoster list all staff with pagination support.
func (env Env) StaffRoster(page int, rowCount int) ([]staff.Account, error) {
	offset := (page - 1) * rowCount
	query := `
	SELECT id AS id,
		username AS userName,
		display_name AS displayName,
		department AS department,
		group_memberships AS groupMembers
	FROM backyard.staff
	WHERE is_active = 1
	ORDER BY id ASC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(query, rowCount, offset)

	var accounts []staff.Account

	if err != nil {
		adminLogger.
			WithField("location", "Query staff roster").
			Error(err)

		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var a staff.Account

		err := rows.Scan(
			&a.ID,
			&a.UserName,
			&a.DisplayName,
			&a.Department,
			&a.GroupMembers,
		)

		if err != nil {
			adminLogger.
				WithField("location", "Staff roster").
				Error(err)

			continue
		}

		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		adminLogger.
			WithField("location", "Staff roster iteration").
			Error(err)

		return accounts, err
	}

	return accounts, nil
}

// ActivateStaff reuses a previously removed staff account
func (env Env) ActivateStaff(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 1
    WHERE username = ?
      AND is_active = 0
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		adminLogger.
			WithField("location", "Activate a staff").
			Error(err)

		return err
	}

	return nil
}

// UpdateStaff updates a staff's profile by administrator
func (env Env) UpdateStaff(userName string, a staff.Account) error {
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
		a.UserName,
		a.Email,
		a.DisplayName,
		a.Department,
		a.GroupMembers,
		userName,
	)

	if err != nil {
		adminLogger.
			WithField("location", "Update staff profile").
			Error(err)

		return err
	}

	return nil
}

// RemoveStaff deactivates a staff's account.
// This is not a SQL DELETE operation.
// It just flags the account as not active.
// When doing this, you should also remove:
// 1. VIP status of all ftc accouts associated with this staff
// 2. All access tokens created by this staff to access next-api
// 3. All access tokens created by this staff to access backyard-api
func (env Env) RemoveStaff(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 0
    WHERE userName = ?
      AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		adminLogger.
			WithField("location", "Deactivate a staff").
			Error(err)

		return err
	}

	return nil
}

// RevokeStaffVIP set vip to false for all ftc accounts associated with a staff
// This should be perfomed when you remove a staff's account.
func (env Env) RevokeStaffVIP(userName string) error {
	query := `
	UPDATE backyard.staff_myft AS s
		LEFT JOIN cmstmp01.userinfo AS u
		ON s.myft_id = u.user_id
	SET isvip = 0
	WHERE s.staff_name = ?
		AND u.isvip = 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		adminLogger.WithField("location", "remove vip status of a staff").Error(err)

		return nil
	}

	return nil
}
