package readers

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env handles FTC user data.
type Env struct {
	db *sqlx.DB
}

func NewEnv(db *sqlx.DB) Env {
	return Env{db: db}
}

func (env Env) BeginMemberTx() (MemberTx, error) {
	tx, err := env.db.Beginx()

	if err != nil {
		return MemberTx{}, err
	}

	return NewMemberTx(tx), nil
}

var logger = logrus.WithField("package", "repository.customer")
