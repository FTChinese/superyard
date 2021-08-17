package admin

import (
	"github.com/FTChinese/superyard/pkg/db"
)

type Env struct {
	dbs db.ReadWriteMyDBs
}

func NewEnv(myDBs db.ReadWriteMyDBs) Env {
	return Env{
		dbs: myDBs,
	}
}
