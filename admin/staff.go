package admin

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/backyard-api/staff"
)

// Create a new staff and generate a random password.
// The password is returned so that you could send it to user's email.
func (env Env) createStaff(a staff.Account, password string) error {

	query := `
	INSERT INTO backyard.staff
      SET username = ?,
        email = ?,
        password = UNHEX(MD5(?)),
        display_name = NULLIF(?, ''),
        department = NULLIF(?, ''),
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
		adminLogger.
			WithField("location", "Inserting new staff").
			Error(err)

		return err
	}

	return nil
}

// NewStaff creates a new account for a staff
// After the account is created, you should send the password to this its email address.
func (env Env) NewStaff(a staff.Account) (postoffice.Parcel, error) {
	password, err := gorest.RandomHex(4)

	if err != nil {
		adminLogger.WithField("location", "Creating password for new staff").Error(err)

		return postoffice.Parcel{}, err
	}

	err = env.createStaff(a, password)

	if err != nil {
		return postoffice.Parcel{}, err
	}

	return a.SignupParcel(password)
}

// StaffRoster list all staff with pagination support.
// Pay attention to SQL nullable columns.
// This API do not provide JSON null to reduce efforts of converting between weak type and Golang's strong type.
// Simply user each type's zero value for JSON nullable fields.
func (env Env) StaffRoster(page int64, rowCount int64) ([]staff.Profile, error) {
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
		adminLogger.
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
			adminLogger.
				WithField("location", "Staff roster").
				Error(err)

			continue
		}

		profiles = append(profiles, p)
	}

	if err := rows.Err(); err != nil {
		adminLogger.
			WithField("location", "Staff roster iteration").
			Error(err)

		return profiles, err
	}

	return profiles, nil
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

// deactivateStaff tuens `is_active` column to false
func (env Env) deactivateStaff(userName string) error {
	query := `
    UPDATE backyard.staff
	  SET is_active = 0,
	  	deactivated_utc = UTC_TIMESTAMP()
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

// revokeStaffVIP set `isvip` column of `userinfo` table to false for all ftc accounts associated with a staff.
// This works only if the staff already associated backyard account with ftc accounts
func (env Env) revokeStaffVIP(userName string) error {
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

// unlinkMyft removes link between CMS account and ftc accounts.
// This is similar to staff.DeleteMyft() in staff/myft.go.
func (env Env) unlinkMyft(userName string) error {
	query := `
	DELETE FROM backyard.staff_myft
    WHERE staff_name = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, userName)

	if err != nil {
		adminLogger.WithField("location", "Deleting staff_myft record").Error(err)
		return err
	}

	return nil
}

// RemoveStaff deactivates a staff's account and optionally revoke VIP status from all ftc accounts associated with this staff
// This is not a SQL DELETE operation.
// It flags the account as not active.
func (env Env) RemoveStaff(userName string, rmVIP bool) error {
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
