package search

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env wraps db for search operations.
type Env struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository.search")
