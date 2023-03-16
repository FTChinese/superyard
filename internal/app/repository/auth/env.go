package auth

import (
	"github.com/FTChinese/superyard/pkg/db"
)

type Env struct {
	DBs db.ReadWriteMyDBs
}

func NewEnv(dbs db.ReadWriteMyDBs) Env {
	return Env{
		DBs: dbs,
	}
}
