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
