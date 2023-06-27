package readers

import (
	"github.com/FTChinese/superyard/pkg/db"
	"go.uber.org/zap"
)

// Env handles FTC user data.
type Env struct {
	gormDBs db.MultiGormDBs
	logger  *zap.Logger
}

func New(myDBs db.MultiGormDBs, l *zap.Logger) Env {
	return Env{
		gormDBs: myDBs,
		logger:  l,
	}
}
