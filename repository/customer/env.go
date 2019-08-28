package customer

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env handles FTC user data.
type Env struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository.customer")
