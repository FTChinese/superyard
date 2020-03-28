package test

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/brianvoe/gofakeit/v4"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/employee"
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
		ID:           employee.GenStaffID(),
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

func (s Staff) Account() employee.Account {
	return employee.Account{
		ID: null.StringFrom(s.ID),
		BaseAccount: employee.BaseAccount{
			UserName:     s.UserName,
			Email:        s.Email,
			DisplayName:  null.StringFrom(s.DisplayName),
			Department:   null.StringFrom(s.Department),
			GroupMembers: s.GroupMembers,
		},
		IsActive: s.IsActive,
	}
}

func (s Staff) Login() employee.Login {
	return employee.Login{
		UserName: s.UserName,
		Password: s.Password,
	}
}

func (s Staff) SignUp() employee.SignUp {
	return employee.SignUp{
		Account:  s.Account(),
		Password: s.Password,
	}
}

func (s Staff) PasswordReset() employee.PasswordReset {
	return employee.PasswordReset{
		Email:    s.Email,
		Token:    s.PwResetToken,
		Password: "",
	}
}

func (s Staff) NewPassword() employee.Credentials {
	return employee.Credentials{
		ID: s.ID,
		Login: employee.Login{
			UserName: s.UserName,
			Password: SimplePassword(),
		},
	}
}

func (s Staff) OldPassword() employee.Credentials {
	return employee.Credentials{
		ID: s.ID,
		Login: employee.Login{
			Password: s.Password,
		},
	}
}
