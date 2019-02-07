package model

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("package", "model")



// Env interact with user data
type StaffEnv struct {
	DB *sql.DB
}