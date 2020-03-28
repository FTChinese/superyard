package apps

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AndroidEnv struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository/apps")
