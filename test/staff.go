//go:build !production

package test

import (
	"time"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/faker"
	"github.com/brianvoe/gofakeit/v5"
)

type Staff struct {
	UserName     string `db:"user_name"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	IsActive     bool
	DisplayName  string `db:"display_name"`
	Department   string `db:"department"`
	GroupMembers int64  `db:"group_memberships"`
	IP           string
	PwResetToken string
}

func NewStaff() Staff {
	gofakeit.Seed(time.Now().UnixNano())
	t, _ := gorest.RandomHex(32)

	return Staff{
		UserName:     gofakeit.Username(),
		Email:        gofakeit.Email(),
		Password:     faker.SimplePassword(),
		IsActive:     true,
		DisplayName:  gofakeit.Name(),
		Department:   "tech",
		GroupMembers: 2,
		IP:           gofakeit.IPv4Address(),
		PwResetToken: t,
	}
}

var FixedStaff = Staff{
	UserName:     "weiguo.ni",
	Email:        "victor@example.org",
	Password:     "12345678",
	IsActive:     false,
	DisplayName:  "Victor",
	Department:   "tech",
	GroupMembers: 2,
	IP:           "127.0.0.1",
	PwResetToken: "",
}
