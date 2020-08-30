package readers

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env handles FTC user data.
type Env struct {
	DB *sqlx.DB
}

func NewEnv(db *sqlx.DB) Env {
	return Env{DB: db}
}

var logger = logrus.WithField("package", "repository.customer")
