package adminmodel

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/staffmodel"
)

// Admin let administrators to manage staff
type Admin struct {
	DB *sql.DB
}

var adminLogger = log.WithFields(log.Fields{
	"package":  "adminmodel",
	"resource": "Admin",
})

// CreateStaff adds create a new staff's profile
func (a Admin) CreateStaff(s Staff) error {
	query := `
	INSERT INTO backyard.staff
      SET username = ?,
        email = ?,
        password = UNHEX(MD5(?)),
        display_name = NULLIF(?, ''),
        department = NULLIF(?, ''),
		group_memberships = ?`

	_, err := a.DB.Exec(query,
		s.UserName,
		s.Email,
		s.Password,
		s.DisplayName,
		s.Department, s.GroupMembers)

	if err != nil {
		adminLogger.
			WithField("func", "CreateStaff").
			Error(err)

		return err
	}

	return nil
}

// ActivateStaff turns a staff as active if it is removed
func (a Admin) ActivateStaff(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 1
    WHERE username = ?
      AND is_active = 0
	LIMIT 1`

	_, err := a.DB.Exec(query, userName)

	if err != nil {
		adminLogger.
			WithField("func", "ActivateStaff").
			Error(err)

		return err
	}

	return nil
}

// DeactivateStaff removes a staff
func (a Admin) DeactivateStaff(userName string) error {
	query := `
    UPDATE backyard.staff
      SET is_active = 0
    WHERE userName = ?
      AND is_active = 1
	LIMIT 1`

	_, err := a.DB.Exec(query, userName)

	if err != nil {
		adminLogger.
			WithField("func", "DeactivateStaff").
			Error(err)

		return err
	}

	return nil
}

// StaffRoster list all staff
func (a Admin) StaffRoster(page int, rowCount int) ([]staffmodel.Account, error) {
	offset := (page - 1) * rowCount
	query := `
	SELECT id AS id,
		username AS userName,
		display_name AS displayName,
		department AS department,
		group_memberships AS groupMembers,
		myft_id AS myftId
	FROM backyard.staff
	WHERE is_active = 1
	ORDER BY id ASC
	LIMIT ? OFFSET ?`

	rows, err := a.DB.Query(query, rowCount, offset)

	var items []staffmodel.Account

	if err != nil {
		adminLogger.
			WithField("func", "StaffRoster").
			Error(err)

		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item staffmodel.Account

		err := rows.Scan(
			&item.ID,
			&item.UserName,
			&item.DisplayName,
			&item.Department,
			&item.Groups,
			&item.MyftID,
		)

		if err != nil {
			adminLogger.
				WithField("func", "StaffRoster").
				Error(err)

			continue
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		adminLogger.
			WithField("func", "StaffRoster").
			Error(err)

		return items, err
	}

	return items, nil
}

// StaffProfile retrieves all data of a staff
func (a Admin) StaffProfile(userName: string) (staff) {

}
