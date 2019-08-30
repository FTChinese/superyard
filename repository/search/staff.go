package search

import (
	"errors"
	"gitlab.com/ftchinese/backyard-api/models/builder"
	"gitlab.com/ftchinese/backyard-api/models/employee"
)

// Staff searches a staff by either email or user name
// ?email=<>
// ?name=<>
func (env Env) Staff(where *builder.Where) (employee.Account, error) {
	if where == nil {
		return employee.Account{}, errors.New("where clause is empty")
	}

	s := sqlSearchStaff + where.Build()

	var account employee.Account
	if err := env.DB.Get(&account, s, where.Values...); err != nil {
		return employee.Account{}, err
	}

	return account, nil
}
