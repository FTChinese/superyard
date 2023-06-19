//go:build !production

package db

import "github.com/FTChinese/superyard/faker"

func MockMySQL() ReadWriteMyDBs {
	faker.MustSetupViper()
	return MustNewMyDBs(false)
}

func MockGormSQL() MultiGormDBs {
	faker.MustSetupViper()
	return MustNewMultiGormDBs(false)
}
