package test

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
)

var DBX *sqlx.DB
var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var CFG = config.Config{Debug: true}

func init() {

	faker.MustConfigViper()

	DBX = db.MustNewDB(CFG.MustGetDBConn(""))
}
