package search

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/repository/staff"
)

// Staff searches a staff by either email or user name
// ?email=<>
// ?name=<>
func (env Env) Staff(col employee.Column, val string) (employee.Account, error) {

	var account employee.Account
	if err := env.DB.Get(&account, staff.QueryAccount(col), val); err != nil {
		return employee.Account{}, err
	}

	return account, nil
}
