package staff

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("package", "staff")

type sqlCol string

const (
	colUserName sqlCol = "username"
	colEmail    sqlCol = "email"
	// This is used by both user login and finding an account
	stmtAccount string = `
	SELECT id AS id,
		username AS userName,
		IFNULL(email, '') AS email,
		IFNULL(display_name, '') AS displayName,
		IFNULL(department, '') AS department,
		group_memberships AS groups
	FROM backyard.staff`
)

const (
	resetLetterURL = "http://localhost:8900/backyard/password-reset"
)

// Env interact with user data
type Env struct {
	DB *sql.DB
}

func (env Env) exists(col sqlCol, value string) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS(
		SELECT *
		FROM backyard.staff
		WHERE %s = ?
	) AS alreadyExists`, string(col))

	var exists bool

	err := env.DB.QueryRow(query, value).Scan(&exists)

	if err != nil {
		logger.
			WithField("location", "staff exists").
			Error(err)

		return false, err
	}

	return exists, nil
}

// StaffNameExists checks if name exists in the username column of backyard.staff table.
func (env Env) StaffNameExists(name string) (bool, error) {
	return env.exists(colUserName, name)
}

// StaffEmailExists checks if an email address exists in the email column of backyard.staff table.
func (env Env) StaffEmailExists(email string) (bool, error) {
	return env.exists(colEmail, email)
}
