package test

import (
	"github.com/brianvoe/gofakeit/v4"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/employee"
)

type Staff struct {
	ID           string
	UserName     string
	Email        string
	Password     string
	IsActive     bool
	DisplayName  string
	Department   string
	GroupMembers int64
}

func NewStaff() Staff {
	return Staff{
		ID:           employee.GenerateID(),
		UserName:     gofakeit.Username(),
		Email:        gofakeit.Email(),
		Password:     "12345678",
		IsActive:     true,
		DisplayName:  gofakeit.Name(),
		Department:   "tech",
		GroupMembers: 2,
	}
}

func (s Staff) Account() employee.Account {
	return employee.Account{
		ID:           null.StringFrom(s.ID),
		Email:        s.Email,
		UserName:     s.UserName,
		Password:     s.Password,
		IsActive:     s.IsActive,
		DisplayName:  null.StringFrom(s.DisplayName),
		Department:   null.StringFrom(s.Department),
		GroupMembers: s.GroupMembers,
	}
}

func (s Staff) Login() employee.Login {
	return employee.Login{
		UserName: s.UserName,
		Password: s.Password,
	}
}
