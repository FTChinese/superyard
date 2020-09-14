package user

import (
	"github.com/jmoiron/sqlx"
)

type Env struct {
	DB *sqlx.DB
}

func NewEnv(db *sqlx.DB) Env {
	return Env{
		DB: db,
	}
}
