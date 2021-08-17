package test

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/faker"
	oauth2 "github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg/staff"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"time"
)

type Staff struct {
	ID           string `db:"staff_id"`
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
		ID:           staff.GenStaffID(),
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
	ID:           "stf_7481cc038eedce2f",
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

func (s Staff) Account() staff.Account {
	return staff.Account{
		ID:           null.StringFrom(s.ID),
		UserName:     s.UserName,
		Email:        s.Email,
		DisplayName:  null.StringFrom(s.DisplayName),
		Department:   null.StringFrom(s.Department),
		GroupMembers: s.GroupMembers,
		IsActive:     s.IsActive,
	}
}

func (s Staff) Credentials() staff.Credentials {
	return staff.Credentials{
		UserName: s.UserName,
		Password: s.Password,
	}
}

func (s Staff) SignUp() staff.SignUp {
	return staff.SignUp{
		Account:  s.Account(),
		Password: s.Password,
	}
}

func (s Staff) PwResetSession() staff.PwResetSession {
	return staff.PwResetSession{
		Email:      s.Email,
		Token:      s.PwResetToken,
		IsUsed:     false,
		ExpiresIn:  10800,
		CreatedUTC: chrono.TimeNow(),
		SourceURL:  "http://localhost:4200/password-reset",
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
