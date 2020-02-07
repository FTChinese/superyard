package staff

import (
	"fmt"
	"gitlab.com/ftchinese/superyard/models/employee"
)

const (
	sqlSelectStaff = `
	SELECT staff_id,
		IFNULL(email, '') AS email,
		user_name,
		is_active,
		display_name,
		department,
		group_memberships`

	stmtLogin = sqlSelectStaff + `
	FROM backyard.staff
	WHERE (user_name, password) = (?, UNHEX(MD5(?)))
		AND is_active = 1`

	// Verify password for a logged-in staff.
	stmtVerifyPassword = sqlSelectStaff + `
	FROM backyard.staff
	WHERE (staff_id, password) = (?, UNHEX(MD5(?)))
		AND is_active = 1`

	// Create a new staff.
	stmtCreateAccount = `
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

	stmtAddID = `
	UPDATE backyard.staff
	SET staff_id = :staff_id
	WHERE user_name = :user_name
	LIMIT 1`

	stmtUpdateLastLogin = `
    UPDATE backyard.staff
      SET last_login_utc = UTC_TIMESTAMP(),
        last_login_ip = IFNULL(INET6_ATON(?), last_login_ip)
    WHERE user_name = ?
	LIMIT 1`

	// Profile
	sqlSelectProfile = sqlSelectStaff + `,
	    created_utc AS created_at,
		deactivated_utc AS deactivated_at,
		updated_utc AS updated_at,
		last_login_utc AS last_login_at,
		INET6_NTOA(staff.last_login_ip) AS last_login_ip
  	FROM backyard.staff`

	// Select a user profile either by user_name,
	// email, or staff_id
	stmtSelectProfile = sqlSelectProfile + `
	WHERE staff_id = ?
	LIMIT 1`

	stmtListStaff = sqlSelectProfile + `
	ORDER BY user_name ASC
	LIMIT ? OFFSET ?`

	// Statement to deactivate a staff
	stmtDeactivate = `
    UPDATE backyard.staff
	  SET is_active = 0,
	  	deactivated_utc = UTC_TIMESTAMP()
    WHERE staff_id = ?
      AND is_active = 1
	LIMIT 1`

	stmtDeletePersonalToken = `
	UPDATE oauth.access
		SET is_active = 0
	WHERE created_by = ?`

	// Restore a deactivated staff.
	stmtActivate = `
    UPDATE backyard.staff
      SET is_active = 1,
      	updated_utc = UTC_TIMESTAMP()
    WHERE staff_id = ?
      AND is_active = 0
	LIMIT 1`

	// Manipulate staff' profile.
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

	stmtUpdatePassword = `
	UPDATE backyard.staff
		SET password = UNHEX(MD5(:password)),
			updated_utc = UTC_TIMESTAMP()
	WHERE staff_id = :staff_id
		AND is_active = 1
	LIMIT 1`

	stmtUpdateLegacyPassword = `
	UPDATE cmstmp01.managers
		SET password = MD5(:password)
	WHERE username = :user_name
	LIMIT 1`

	// Password reset token
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
)

// QueryAccount build SQL to retrieve a staff's account
// by one of staff_id, email or user_name columns.
func QueryAccount(col employee.Column) string {
	return fmt.Sprintf(sqlSelectStaff+`
	FROM backyard.staff
	WHERE %s = ?
	LIMIT 1`, col.String())
}
