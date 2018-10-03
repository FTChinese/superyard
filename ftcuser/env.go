package ftcuser

import (
	"database/sql"
)

// Env wraps DB connection
type Env struct {
	DB *sql.DB
}
