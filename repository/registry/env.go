package registry

import (
	"github.com/jmoiron/sqlx"
)

// Env wraps db.
type Env struct {
	DB *sqlx.DB
}
