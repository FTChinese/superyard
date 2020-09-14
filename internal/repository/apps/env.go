package apps

import (
	"github.com/jmoiron/sqlx"
)

type AndroidEnv struct {
	DB *sqlx.DB
}
