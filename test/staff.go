//go:build !production

package test

import (
	"time"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/faker"
	oauth2 "github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
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

func (s Staff) Account() user.Account {
	return user.Account{
		UserName:    s.UserName,
		Email:       s.Email,
		DisplayName: s.DisplayName,
	}
}

func (s Staff) PwResetSession() user.PwResetSession {
	return user.PwResetSession{
		Email:      s.Email,
		Token:      conv.HexStr(s.PwResetToken),
		IsUsed:     false,
		ExpiresIn:  10800,
		CreatedUTC: chrono.TimeNow(),
	}
}

func (s Staff) MustNewOAuthApp() oauth2.App {

	app, err := oauth2.NewApp(genOAuthApp(), s.UserName)

	if err != nil {
		panic(err)
	}

	return app
}

func (s Staff) MustNewPersonalKey() oauth2.Access {
	key, err := oauth2.NewAccess(oauth2.BaseAccess{
		Description: null.StringFrom(gofakeit.Sentence(10)),
		ClientID:    null.String{},
	}, s.UserName)

	if err != nil {
		panic(err)
	}

	return key
}

func (s Staff) MustNewAppToken(app oauth2.App) oauth2.Access {
	token, err := oauth2.NewAccess(oauth2.BaseAccess{
		Description: null.StringFrom(gofakeit.Sentence(10)),
		ClientID:    null.StringFrom(app.ClientID),
	}, s.UserName)

	if err != nil {
		panic(err)
	}

	return token
}
