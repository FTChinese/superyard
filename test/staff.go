package test

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/guregu/null"
	"github.com/icrowley/fake"
	"gitlab.com/ftchinese/backyard-api/models/employee"
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
	id, err := gorest.RandomHex(8)
	if err != nil {
		panic(err)
	}

	return Staff{
		ID:           id,
		UserName:     fake.UserName(),
		Email:        fake.EmailAddress(),
		Password:     "12345678",
		IsActive:     true,
		DisplayName:  fake.UserName(),
		Department:   "tech",
		GroupMembers: 2,
	}
}

func (s Staff) Account() employee.Account {
	return employee.Account{
		ID:           s.ID,
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