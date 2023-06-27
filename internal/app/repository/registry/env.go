package registry

import (
	"github.com/FTChinese/superyard/pkg/db"
)

// Env wraps db.
type Env struct {
	gormDBs db.MultiGormDBs
}

func NewEnv(myDBs db.MultiGormDBs) Env {
	return Env{
		gormDBs: myDBs,
	}
}
