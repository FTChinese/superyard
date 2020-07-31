package test

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/oauth"
	"gitlab.com/ftchinese/superyard/pkg/staff"
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
		Password:     simplePassword(),
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

func (s Staff) MustNewOAuthApp() oauth.App {

	app, err := oauth.NewApp(genOAuthApp(), s.UserName)

	if err != nil {
		panic(err)
	}

	return app
}

func (s Staff) MustNewPersonalKey() oauth.Access {
	key, err := oauth.NewAccess(oauth.BaseAccess{
		Description: null.StringFrom(gofakeit.Sentence(10)),
		ClientID:    null.String{},
	}, s.UserName)

	if err != nil {
		panic(err)
	}

	return key
}

func (s Staff) MustNewAppToken(app oauth.App) oauth.Access {
	token, err := oauth.NewAccess(oauth.BaseAccess{
		Description: null.StringFrom(gofakeit.Sentence(10)),
		ClientID:    null.StringFrom(app.ClientID),
	}, s.UserName)

	if err != nil {
		panic(err)
	}

	return token
}
