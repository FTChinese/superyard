package registry

import (
	"github.com/FTChinese/superyard/pkg/db"
)

// Env wraps db.
type Env struct {
	dbs db.ReadWriteMyDBs
}

func NewEnv(myDBs db.ReadWriteMyDBs) Env {
	return Env{
		dbs: myDBs,
	}
}