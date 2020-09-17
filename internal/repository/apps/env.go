package apps

import (
	"github.com/jmoiron/sqlx"
)

type Env struct {
	DB *sqlx.DB
}
