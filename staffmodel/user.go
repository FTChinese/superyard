package staffmodel

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

var staffLogger = log.WithFields(log.Fields{
	"package":  "staffmodel",
	"resource": "User",
})

// User interact with user data
type User struct {
	DB *sql.DB
}

// Auth perform authentication by user name and password
func (u User) Auth(l Login) (Account, error) {
	query := `
	SELECT id AS id,
		username AS userName,
		display_name AS displayName,
		department AS department,
		group_memberships AS groups,
		vip_uuid AS myftId
	FROM backyard.staff
	WHERE username = ?
		AND password = UNHEX(MD5(?))
		AND is_active = 1
	LIMIT 1`

	var a Account
	err := u.DB.QueryRow(query, l.UserName, l.Password).Scan(
		&a.ID,
		&a.UserName,
		&a.DisplayName,
		&a.Department,
		&a.Groups,
		&a.MyftID,
	)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "Auth",
			"table": "backyard.staff",
		}).Error(err)

		return a, err
	}

	go u.UpdateLoginHistory(l)

	return a, nil
}

// UpdateLoginHistory saves user login footprint after successfully authenticated.
func (u User) UpdateLoginHistory(l Login) error {
	query := `
    UPDATE backyard.staff
      SET last_login_utc = UTC_TIMESTAMP(),
        last_login_ip = IFNULL(INET6_ATON(:?), last_login_ip)
    WHERE username = :?
	LIMIT 1`

	_, err := u.DB.Exec(query, l.UserIP, l.UserName)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdateLoginHistory",
			"table": "backyard.staff",
		}).Error(err)

		return err
	}

	return nil
}

// UpdatePassword allows a user to change password after login
func (u User) UpdatePassword(p Password) error {
	query := `
	UPDATE backyard.staff
		SET password = UNHEX(MD5(?)),
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
		AND is_active = 1
	LIMIT 1`

	_, err := u.DB.Exec(query, p.New, p.UserName)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdatePassword",
			"table": "backyard.staff",
		}).Error(err)
		return err
	}

	legacyQuery := `
	UPDATE cmstmp01.managers
		SET password = MD5(?)
	WHERE username = ?
	LIMIT 1`

	_, err = u.DB.Exec(legacyQuery, p.New, p.UserName)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdatePassword",
			"table": "cmstmp01.managers",
		}).Error(err)
		return err
	}

	return nil
}

// UpdateName allows a user to change its display name.
// Not userName cannot be changed.
func (u User) UpdateName(p Profile) error {
	query := `
	UPDATE backyard.staff
		SET display_name = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
	LIMIT 1`

	_, err := u.DB.Exec(query, p.DisplayName, p.UserName)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdateName",
			"table": "backyard.staff",
		}).Error(err)
		return err
	}

	return nil
}

// UpdateEmail allows a user to udpate its email address.
func (u User) UpdateEmail(p Profile) error {
	query := `
	UPDATE backyard.staff
		SET email = ?
			updated_utc = UTC_TIMESTAMP()
	WHERE username = ?
	LIMIT 1`

	_, err := u.DB.Exec(query, p.DisplayName, p.UserName)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "UpdateEmail",
			"table": "backyard.staff",
		}).Error(err)
		return err
	}

	return nil
}

// Profile retrieves all of a user's data
func (u User) Profile(userName string) (Profile, error) {
	query := `
	SELECT staff.id AS id,
		staff.username AS userName,
		staff.email AS email,
		staff.is_active AS isActive,
		staff.display_name AS displayName,
		staff.department AS department,
		staff.group_memberships AS groupMembers,
		staff.vip_uuid AS vipId,
		myft.email AS vipEmail,
		staff.created_utc AS createdAt,
		staff.deactivated_utc AS deactivatedAt,
		staff.updated_utc AS updatedAt,
		staff.last_login_utc AS lastLoginAt,
		INET6_NTOA(staff.last_login_ip) AS lastLoginIp
  	FROM backyard.staff AS staff
    	LEFT JOIN cmstmp01.userinfo AS myft
		ON staff.vip_uuid = myft.user_id
	WHERE username = ?
	LIMIT 1`

	var p Profile
	err := u.DB.QueryRow(query, userName).Scan(
		&p.ID,
		&p.UserName,
		&p.Email,
		&p.IsActive,
		&p.DisplayName,
		&p.Department,
		&p.GroupMembers,
		&p.MyftID,
		&p.MyftEmail,
		&p.CreatedAt,
		&p.DeactiviateAt,
		&p.UpdatedAt,
		&p.LastLoginAt,
		&p.LastLoginIP,
	)

	if err != nil {
		staffLogger.WithFields(log.Fields{
			"func":  "Profile",
			"table": "backyard.staff",
		}).Error(err)

		return p, err
	}

	return p, nil
}

// ResetPassword allows a user to reset password in case of forgotten.
func (u User) ResetPassword() {

}
