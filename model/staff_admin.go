package model

import "gitlab.com/ftchinese/backyard-api/staff"

// StaffRoster list all staff with pagination support.
// Pay attention to SQL nullable columns.
// This API do not provide JSON null to reduce efforts of converting between weak type and Golang's strong type.
// Simply user each type's zero value for JSON nullable fields.
func (env StaffEnv) StaffRoster(page int64, rowCount int64) ([]staff.Profile, error) {
	offset := (page - 1) * rowCount
	query := `
	SELECT id AS id,
		username AS userName,
		IFNULL(email, '') AS email,
		IFNULL(display_name, '') AS displayName,
		IFNULL(department, '') AS department,
		group_memberships AS groupMembers,
		is_active AS isActive
	FROM backyard.staff
	ORDER BY id ASC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(query, rowCount, offset)

	if err != nil {
		logger.
			WithField("location", "Query staff roster").
			Error(err)

		return nil, err
	}
	defer rows.Close()

	profiles := make([]staff.Profile, 0)
	for rows.Next() {
		var p staff.Profile

		err := rows.Scan(
			&p.ID,
			&p.UserName,
			&p.Email,
			&p.DisplayName,
			&p.Department,
			&p.GroupMembers,
			&p.IsActive,
		)

		if err != nil {
			logger.
				WithField("location", "Staff roster").
				Error(err)

			continue
		}

		profiles = append(profiles, p)
	}

	if err := rows.Err(); err != nil {
		logger.
			WithField("location", "Staff roster iteration").
			Error(err)

		return profiles, err
	}

	return profiles, nil
}

// UpdateStaff updates a staff's profile by administrator
func (env StaffEnv) UpdateStaff(userName string, a staff.Account) error {
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
		logger.
			WithField("location", "Update staff profile").
			Error(err)

		return err
	}

	return nil
}

// deactivateStaff tuens `is_active` column to false
func (env StaffEnv) deactivateStaff(userName string) error {
	query := `
    UPDATE backyard.staff
	  SET is_active = 0,
	  	deactivated_utc = UTC_TIMESTAMP()
    WHERE userName = ?
      AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		logger.
			WithField("location", "Deactivate a staff").
			Error(err)

		return err
	}

	return nil
}

// revokeStaffVIP set `isvip` column of `userinfo` table to false for all ftc accounts associated with a staff.
// This works only if the staff already associated backyard account with ftc accounts
func (env StaffEnv) revokeStaffVIP(userName string) error {
	query := `
	UPDATE backyard.staff_myft AS s
		LEFT JOIN cmstmp01.userinfo AS u
		ON s.myft_id = u.user_id
	SET isvip = 0
	WHERE s.staff_name = ?
		AND u.isvip = 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		logger.WithField("location", "remove vip status of a staff").Error(err)

		return nil
	}

	return nil
}

// unlinkMyft removes link between CMS account and ftc accounts.
// This is similar to staff.DeleteMyft() in staff/myft.go.
func (env StaffEnv) unlinkMyft(userName string) error {
	query := `
	DELETE FROM backyard.staff_myft
    WHERE staff_name = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		logger.WithField("location", "Deleting staff_myft record").Error(err)
		return err
	}

	return nil
}

// RemoveStaff deactivates a staff's account and optionally revoke VIP status from all ftc accounts associated with this staff
// This is not a SQL DELETE operation.
// It flags the account as not active.
func (env StaffEnv) RemoveStaff(userName string, rmVIP bool) error {
	if rmVIP {
		err := env.revokeStaffVIP(userName)

		if err != nil {
			return err
		}
	}

	err := env.unlinkMyft(userName)
	if err != nil {
		return err
	}

	err = env.deactivateStaff(userName)

	if err != nil {
		return err
	}

	return nil
}

// ActivateStaff reuses a previously removed staff account
func (env StaffEnv) ActivateStaff(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 1
    WHERE username = ?
      AND is_active = 0
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		logger.
			WithField("location", "Activate a staff").
			Error(err)

		return err
	}

	return nil
}

// VIPRoster list all vip account on ftchinese.com
func (env StaffEnv) VIPRoster() ([]staff.MyftAccount, error) {
	query := `
	SELECT user_id AS id,
		email AS email
	FROM cmstmp01.userinfo
	WHERE is_vip = 1`

	rows, err := env.DB.Query(query)

	if err != nil {
		logger.WithField("location", "Query myft vip accounts").Error(err)

		return nil, err
	}
	defer rows.Close()

	vips := make([]staff.MyftAccount, 0)
	for rows.Next() {
		var vip staff.MyftAccount

		err := rows.Scan(
			&vip.ID,
			&vip.Email,
		)

		if err != nil {
			logger.WithField("location", "Scan myft vip account").Error(err)

			continue
		}

		vips = append(vips, vip)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "rows iteration").Error(err)

		return vips, err
	}

	return vips, nil
}

func (env StaffEnv) updateVIP(myftID string, isVIP bool) error {
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
func (env StaffEnv) GrantVIP(myftID string) error {
	return env.updateVIP(myftID, true)
}

// RevokeVIP removes vip status from a ftc account
func (env StaffEnv) RevokeVIP(myftID string) error {
	return env.updateVIP(myftID, false)
}