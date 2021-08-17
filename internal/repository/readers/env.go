package readers

import (
	"github.com/FTChinese/superyard/pkg/db"
	"go.uber.org/zap"
)

// Env handles FTC user data.
type Env struct {
	dbs    db.ReadWriteMyDBs
	logger *zap.Logger
}

func NewEnv(myDBs db.ReadWriteMyDBs, l *zap.Logger) Env {
	return Env{
		dbs:    myDBs,
		logger: l,
	}
}

func (env Env) BeginMemberTx() (MemberTx, error) {
	tx, err := env.dbs.Delete.Beginx()

	if err != nil {
		return MemberTx{}, err
	}

	return NewMemberTx(tx), nil
}
