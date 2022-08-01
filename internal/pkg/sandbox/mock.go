//go:build !production

package sandbox

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
)

func MockTestAccount() TestAccount {
	faker.SeedGoFake()
	return TestAccount{
		FtcID:         uuid.New().String(),
		Email:         gofakeit.Email(),
		ClearPassword: "12345678",
		CreatedBy:     "weiguo.ni",
	}
}
