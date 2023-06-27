//go:build !production

package db

import "github.com/FTChinese/superyard/faker"

func MockGormSQL() MultiGormDBs {
	faker.MustSetupViper()
	return MustNewMultiGormDBs(false)
}
