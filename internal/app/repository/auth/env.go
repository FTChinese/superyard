package auth

import (
	"github.com/FTChinese/superyard/pkg/db"
)

type Env struct {
	gormDBs db.MultiGormDBs
}

func NewEnv(gormDBs db.MultiGormDBs) Env {
	return Env{
		gormDBs: gormDBs,
	}
}
