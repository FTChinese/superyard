package model

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("package", "model")

type sqlCol string

const (
	colStaffName sqlCol = "username"
	colEmail     sqlCol = "email"
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

// Env interact with user data
type StaffEnv struct {
	DB *sql.DB
}