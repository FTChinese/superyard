package staff

import "gitlab.com/ftchinese/backyard-api/util"

// Profile contains the full data of a staff
type Profile struct {
	ID            int    `json:"id"`
	UserName      string `json:"userName"`
	Email         string `json:"email"`
	IsActive      bool   `json:"isActive"`
	DisplayName   string `json:"displayName"`
	Department    string `json:"department"`
	GroupMembers  int    `json:"groupMembers"`
	CreatedAt     string `json:"createdAt"`
	DeactiviateAt string `json:"deactivatedAt"`
	UpdatedAt     string `json:"updatedAt"`
	LastLoginAt   string `json:"lastLoginAt"`
	LastLoginIP   string `json:"lastLoginIp"`
}

// Profile retrieves all of a user's data.
// This is used by both an administrator or the user itself
// GET /user/profile
// GET /staff/profile
func (env Env) Profile(userName string) (Profile, error) {
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

// UpdateName allows a user to change its display name.
// PATCH /user/display-name
func (env Env) UpdateName(userName string, displayName string) error {
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
func (env Env) UpdateEmail(userName string, email string) error {
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

// UpdatePassword allows user to change password in its settings.
func (env Env) UpdatePassword(userName string, p Password) error {
	// Verify user's old password
	matched, err := env.isPasswordMatched(userName, p.Old)

	if err != nil {
		return err
	}

	// Tells controller to respond with 403 Forbidden
	if !matched {
		return util.ErrWrongPassword
	}

	err = env.changePassword(userName, p.New)

	if err != nil {
		return err
	}

	return nil
}
