package staffmodel

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// CMSUser interact with user data
type CMSUser struct {
	DB *sql.DB
}

// Auth perform authentication by user name and password
func (u CMSUser) Auth(l StaffLogin) (StaffAccount, error) {
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

	var a StaffAccount
	err := u.DB.QueryRow(query, l.UserName, l.Password).Scan(
		&a.ID,
		&a.UserName,
		&a.DisplayName,
		&a.Department,
		&a.Groups,
		&a.MyftID,
	)

	if err != nil {
		log.Printf("CMSUser Auth error: %v\n", err)

		return a, err
	}

	return a, nil
}
