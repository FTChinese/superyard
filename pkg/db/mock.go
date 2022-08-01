package db

import "github.com/FTChinese/superyard/faker"

func MockMySQL() ReadWriteMyDBs {
	faker.MustSetupViper()
	return MustNewMyDBs(false)
}
