package ftcuser

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("package", "ftcuser")

// Env wraps DB connection
type Env struct {
	DB *sql.DB
}

type sqlCol string

const (
	colUserName = "user_name"
	colEmail    = "email"
)
