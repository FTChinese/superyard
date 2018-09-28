package adminmodel

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Admin let administrators to manage staff
type Admin struct {
	DB *sql.DB
}

var adminLogger = log.WithFields(log.Fields{
	"package":  "adminmodel",
	"resource": "Admin",
})

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
			WithField("func", "createStaff").
			Error(err)

		return err
	}

	return nil
}

// ActivateStaff turns a staff as active if it is removed
func (a Admin) ActivateStaff() {

}

// DeactivateStaff removes a staff
func (a Admin) DeactivateStaff() {

}

// StaffRoster list all staff
func (a Admin) StaffRoster() {

}

// StaffProfile retrieves all data of a staff
func (a Admin) StaffProfile() {

}
