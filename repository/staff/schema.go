package staff

import "fmt"

const (
	stmtLogin = `
	SELECT staff_id,
		IFNULL(email, '') AS email,
		user_name,
		is_active,
		display_name,
		department,
		group_memberships
	FROM backyard.staff
	WHERE (user_name, password) = (?, UNHEX(MD5(?)))
		AND is_active = 1`

	stmtVerifyPassword = `
	SELECT EXISTS(
		SELECT *
		FROM backyard.staff
		WHERE (user_name, password) = (?, UNHEX(MD5(?))
			AND is_active = 1)`

	stmtUpdateLastLogin = `
    UPDATE backyard.staff
      SET last_login_utc = UTC_TIMESTAMP(),
        last_login_ip = IFNULL(INET6_ATON(?), last_login_ip)
    WHERE user_name = ?
	LIMIT 1`

	stmtInsertResetToken = `
	INSERT INTO backyard.password_reset
    SET token = UNHEX(?),
		email = ?,
		created_utc = UTC_TIMESTAMP()`

	stmtSelectResetToken = `
	SELECT email
	FROM backyard.password_reset
    WHERE token = UNHEX(?)
      AND is_used = 0
	  AND DATE_ADD(created_utc, INTERVAL expires_in SECOND) > UTC_TIMESTAMP()
	  AND is_active = 1
	LIMIT 1`

	stmtDeleteResetToken = `
	UPDATE backyard.password_reset
	SET is_used = 1
    WHERE token = UNHEX(?)
	LIMIT 1`

	// Create a new staff.
	stmtInsertEmployee = `
	INSERT INTO backyard.staff
      SET staff_id = :staff_id,
		user_name = :user_name,
        email = :email,
        password = UNHEX(MD5(:password)),
        display_name = :display_name,
        department = :department,
		group_memberships = :group_memberships,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	stmtSelectProfile = `
	SELECT staff_id,
		IFNULL(email, '') AS email,
		user_name,
		is_active,
		display_name,
		department,
		group_memberships,
	    created_utc AS created_at,
		deactivated_utc AS deactivated_at,
		updated_utc AS updated_at,
		last_login_utc AS last_login_at,
		INET6_NTOA(staff.last_login_ip) AS last_login_ip
  	FROM backyard.staff`

	// Select a user profile either by user_name,
	// email, or staff_id
	stmtIndividualProfile = stmtSelectProfile + `
	WHERE %s = ?
	LIMIT 1`

	stmtListStaff = stmtSelectProfile + `
	ORDER BY user_name ASC
	LIMIT ? OFFSET ?`

	stmtUpdateProfile = `
	UPDATE backyard.staff
	SET user_name = :user_name,
		email = :email,
		display_name = :display_name,
		department = :department,
		group_memberships = :group_memberships,
		updated_utc = UTC_TIMESTAMP()
	WHERE staff_id = :staff_id
		AND is_active = 1
	LIMIT 1`

	stmtDeactivate = `
    UPDATE backyard.staff
	  SET is_active = 0,
	  	deactivated_utc = UTC_TIMESTAMP()
    WHERE staff_id = ?
      AND is_active = 1
	LIMIT 1`

	stmtRevokeVIP = `
	UPDATE backyard.staff_myft AS s
		LEFT JOIN cmstmp01.userinfo AS u
		ON s.myft_id = u.user_id
	SET is_vip = 0
	WHERE s.staff_id = ?`

	stmtDeleteMyft = `
	DELETE FROM backyard.staff_myft
    WHERE staff_id = ?`

	stmtDeletePersonalToken = `
	UPDATE oauth.access
		SET is_active = 0
	WHERE created_by = ?`

	stmtActivate = `
    UPDATE backyard.staff
      SET is_active = 1,
      	updated_utc = UTC_TIMESTAMP()
    WHERE user_name = ?
      AND is_active = 0
	LIMIT 1`

	stmtUpdateName = `
	UPDATE backyard.staff
		SET display_name = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE staff_id = ?
		AND is_active = 1
	LIMIT 1`

	stmtUpdateEmail = `
	UPDATE backyard.staff
		SET email = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE user_name = ?
		AND is_active = 1
	LIMIT 1`

	stmtUpdatePassword = `
	UPDATE backyard.staff
		SET password = UNHEX(MD5(?)),
			updated_utc = UTC_TIMESTAMP()
	WHERE user_name = ?
		AND is_active = 1
	LIMIT 1`

	stmtUpdateLegacyPassword = `
	UPDATE cmstmp01.managers
		SET password = MD5(?)
	WHERE username = ?
	LIMIT 1`

	stmtAuthFtc = `
	SELECT user_id,
		email,
		is_vip
	FROM cmstmp01.userinfo
	WHERE (email, password) = (?, MD5(?))
	LIMIT 1`

	stmtLinkFtc = `
	INSERT INTO backyard.staff_myft
    SET staff_id = :staff_id,
		myft_id = :myft_id,
		created_utc = UTC_TIMESTAMP()`

	stmtSelectFtc = `
	SELECT u.user_id AS user_id,
		u.email AS email,
	    u.is_vip AS is_vip
    FROM backyard.staff_myft AS s
      INNER JOIN cmstmp01.userinfo AS u
      ON s.myft_id = u.user_id
	WHERE s.staff_id = ?`

	stmtDeleteFtc = `
	DELETE FROM backyard.staff_myft
	WHERE staff_id = :staff_id
		AND myft_id = :myft_id
	LIMIT 1`
)

func queryProfile(col Column) string {
	return fmt.Sprintf(stmtIndividualProfile, string(col))
}
