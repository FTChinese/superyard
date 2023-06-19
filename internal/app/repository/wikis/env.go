package wikis

import (
	"github.com/FTChinese/superyard/pkg/db"
)

type Env struct {
	dbs     db.ReadWriteMyDBs
	gormDBs db.MultiGormDBs
}

func NewEnv(myDBs db.ReadWriteMyDBs, gormDBs db.MultiGormDBs) Env {
	return Env{
		dbs:     myDBs,
		gormDBs: gormDBs,
	}
}
