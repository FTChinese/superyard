package auth

import (
	"github.com/FTChinese/superyard/pkg/db"
)

type Env struct {
	DBs     db.ReadWriteMyDBs
	gormDBs db.MultiGormDBs
}

func NewEnv(dbs db.ReadWriteMyDBs, gormDBs db.MultiGormDBs) Env {
	return Env{
		DBs:     dbs,
		gormDBs: gormDBs,
	}
}
