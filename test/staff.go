package test

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/brianvoe/gofakeit/v4"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/staff"
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
		Password:     "12345678",
		IsActive:     true,
		DisplayName:  gofakeit.Name(),
		Department:   "tech",
		GroupMembers: 2,
		IP:           gofakeit.IPv4Address(),
		PwResetToken: t,
	}
}

func (s Staff) Account() staff.Account {
	return staff.Account{
		ID: null.StringFrom(s.ID),
		BaseAccount: staff.BaseAccount{
			UserName:     s.UserName,
			Email:        s.Email,
			DisplayName:  null.StringFrom(s.DisplayName),
			Department:   null.StringFrom(s.Department),
			GroupMembers: s.GroupMembers,
		},
		IsActive: s.IsActive,
	}
}

func (s Staff) Login() staff.Login {
	return staff.Login{
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

func (s Staff) PasswordReset() staff.PasswordReset {
	return staff.PasswordReset{
		Email:    s.Email,
		Token:    s.PwResetToken,
		Password: "",
	}
}

func (s Staff) NewPassword() staff.Credentials {
	return staff.Credentials{
		ID: s.ID,
		Login: staff.Login{
			UserName: s.UserName,
			Password: SimplePassword(),
		},
	}
}

func (s Staff) OldPassword() staff.Credentials {
	return staff.Credentials{
		ID: s.ID,
		Login: staff.Login{
			Password: s.Password,
		},
	}
}
