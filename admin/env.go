package admin

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Env wraps a database connection
type Env struct {
	DB *sql.DB
}

var adminLogger = log.WithFields(log.Fields{
	"package": "adminmodel",
})

const newStaffLetterURL = "http://localhost:8001/backyard/new-staff"

type sqlCol int

const (
	staffNameCol  sqlCol = 0
	staffEmailCol sqlCol = 1
)

func (col sqlCol) String() string {
	cols := [...]string{
		"username",
		"email",
	}

	return cols[col]
}
