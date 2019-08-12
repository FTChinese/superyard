package staff

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Env struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "model.staff")
