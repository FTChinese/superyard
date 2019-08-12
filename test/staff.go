package test

import (
	"github.com/guregu/null"
	"github.com/icrowley/fake"
	"gitlab.com/ftchinese/backyard-api/models/employee"
)

func GenEmployee() employee.Account {
	a, err := employee.NewAccount()
	if err != nil {
		panic(err)
	}

	a.Email = fake.EmailAddress()
	a.UserName = fake.UserName()
	a.Password = "12345678"
	a.IsActive = true
	a.DisplayName = null.StringFrom(fake.UserName())
	a.Department = null.StringFrom("tech")
	a.GroupMembers = 2

	return a
}
