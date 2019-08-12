package admin

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Env struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("trace", "model.admin")
