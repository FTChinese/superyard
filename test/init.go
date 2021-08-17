package test

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/db"
	"math/rand"
	"time"
)

var DBX db.ReadWriteMyDBs
var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {

	faker.MustConfigViper()

	DBX = db.MustNewMyDBs(false)
}
