package model

import (
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
)

// IsPasswordMatched checks whether a staff's credentials are correct.
func (env StaffEnv) IsPasswordMatched(userName, password string) (bool, error) {
	query := `
	SELECT password = UNHEX(MD5(?)) AS matched
	FROM backyard.staff
	WHERE user_name = ?
		AND is_active = 1
	LIMIT 1`

	var matched bool
	err := env.DB.QueryRow(query, password, userName).Scan(&matched)

	if err != nil {
		logger.WithField("trace", "IsPasswordMatched").Error(err)

		return false, err
	}

	return matched, nil
}

// UpdateLoginHistory saves user login footprint after successfully authenticated.
func (env StaffEnv) UpdateLoginHistory(l staff.Login, ip string) error {
	query := `
    UPDATE backyard.staff
      SET last_login_utc = UTC_TIMESTAMP(),
        last_login_ip = IFNULL(INET6_ATON(?), last_login_ip)
    WHERE user_name = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, ip, l.UserName)

	if err != nil {
		logger.WithField("trace", "UpdateLoginHistory").Error(err)
		return err
	}

	return nil
}

// Change password is used by both UpdatePassword after user logged in, or reset password if user forgot it.
func (env StaffEnv) changePassword(userName string, password string) error {
	tx, err := env.DB.Begin()

	query := `
	UPDATE backyard.staff
		SET password = UNHEX(MD5(?)),
			updated_utc = UTC_TIMESTAMP()
	WHERE user_name = ?
		AND is_active = 1
	LIMIT 1`

	_, err = tx.Exec(query, password, userName)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "changePassword").Error(err)
	}

	query = `
	UPDATE cmstmp01.managers
		SET password = MD5(?)
	WHERE username = ?
	LIMIT 1`

	_, err = tx.Exec(query, password, userName)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "changePassword").Error(err)
	}

	if err := tx.Commit(); err != nil {
		logger.WithField("trace", "changePassword").Error(err)
		return err
	}

	return nil
}

// SavePwResetToken send a password reset token to a user's email
func (env StaffEnv) SavePwResetToken(h staff.TokenHolder) error {

	query := `
	INSERT INTO backyard.password_reset
    SET token = UNHEX(?),
		email = ?,
		created_utc = UTC_TIMESTAMP()`

	_, err := env.DB.Exec(query, h.GetToken(), h.GetEmail())

	if err != nil {
		logger.WithField("trace", "SavePwResetToken").Error(err)
		return err
	}

	return nil
}

// VerifyResetToken finds the account associated with a password reset token
// If an account associated with a token is found, allow user to reset password.
func (env StaffEnv) VerifyResetToken(token string) (staff.Account, error) {
	query := `
	SELECT s.id AS id,
	    s.user_name AS userName,
		IFNULL(s.email, '') AS email,
		s.display_name AS displayName,
	    s.department AS department,
	    s.group_memberships
	FROM backyard.password_reset AS t
		LEFT JOIN backyard.staff AS s
		ON t.email = s.email
    WHERE t.token = UNHEX(?)
      AND t.is_used = 0
	  AND DATE_ADD(t.created_utc, INTERVAL t.expires_in SECOND) > UTC_TIMESTAMP()
	  AND s.is_active = 1
	LIMIT 1`

	var a staff.Account
	err := env.DB.QueryRow(query, token).Scan(
		&a.ID,
		&a.UserName,
		&a.Email,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
	)

	if err != nil {
		logger.WithField("trace", "VerifyResetToken").Error(err)

		return a, err
	}

	return a, nil
}

// ResetPassword allows user to reset password after clicked the password reset link in its email.
func (env StaffEnv) ResetPassword(r staff.PasswordReset) error {
	// First try to find the account associated with the token
	a, err := env.VerifyResetToken(r.Token)

	// SqlNoRows
	if err != nil {
		return err
	}

	// The account associated with a token is found. Chnage password.
	err = env.changePassword(a.UserName, r.Password)

	if err != nil {
		return err
	}

	err = env.deleteResetToken(r.Token)

	if err != nil {
		return err
	}

	return nil
}

// DeleteResetToken deletes a password reset token after it was used.
func (env StaffEnv) deleteResetToken(token string) error {
	query := `
	UPDATE backyard.password_reset
	SET is_used = 1
    WHERE token = UNHEX(?)
	LIMIT 1`

	_, err := env.DB.Exec(query, token)

	if err != nil {
		logger.WithField("location", "Deleting a used password reset token").Error(err)

		return err
	}

	return nil
}

// UpdatePassword allows user to change password in its settings.
func (env StaffEnv) UpdatePassword(userName string, p staff.Password) error {
	// Verify user's old password
	matched, err := env.IsPasswordMatched(userName, p.Old)

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
