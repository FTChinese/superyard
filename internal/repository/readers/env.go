package readers

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Env handles FTC user data.
type Env struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewEnv(db *sqlx.DB, l *zap.Logger) Env {
	return Env{
		db:     db,
		logger: l,
	}
}

func (env Env) BeginMemberTx() (MemberTx, error) {
	tx, err := env.db.Beginx()

	if err != nil {
		return MemberTx{}, err
	}

	return NewMemberTx(tx), nil
}
