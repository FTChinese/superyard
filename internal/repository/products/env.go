package products

import (
	"github.com/FTChinese/superyard/pkg/db"
)

type Env struct {
	dbs db.ReadWriteMyDBs
}

func NewEnv(dbs db.ReadWriteMyDBs) Env {
	return Env{
		dbs: dbs,
	}
}
