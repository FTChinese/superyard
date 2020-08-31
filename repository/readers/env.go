package readers

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env handles FTC user data.
type Env struct {
	db *sqlx.DB
}

func NewEnv(db *sqlx.DB) Env {
	return Env{db: db}
}

var logger = logrus.WithField("package", "repository.customer")
