package staff

import (
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

// MyftAuth authenticate a user's myft account.
// Returns staff's FtcAccount if found.
func (env Env) MyftAuth(l reader.Login) (employee.FtcAccount, error) {

	var ftcAccount employee.FtcAccount
	err := env.DB.Get(
		&ftcAccount,
		stmtAuthFtc,
		l.Email,
		l.Password)

	if err != nil {
		logger.WithField("trace", "Env.MyftAuth").Error(err)
		return ftcAccount, err
	}

	return ftcAccount, nil
}

// LinkFtc authenticate a myft account and associated it with a staff account in passed.
func (env Env) LinkFtc(linked employee.Myft) error {
	_, err := env.DB.NamedExec(stmtLinkFtc, &linked)

	if err != nil {
		return err
	}

	return nil
}

// ListMyft lists all myft accounts owned by a staff.
func (env Env) ListMyft(staffID string) ([]employee.FtcAccount, error) {

	accounts := []employee.FtcAccount{}

	err := env.DB.Select(
		&accounts,
		stmtSelectFtc,
		staffID)

	if err != nil {
		logger.WithField("trace", "Env.ListMyft").Error(err)
		return nil, err
	}

	return accounts, nil
}

// UnlinkFtc allows a user to delete a myft account
func (env Env) UnlinkFtc(my employee.Myft) error {

	_, err := env.DB.NamedExec(stmtDeleteFtc, my)

	if err != nil {
		logger.WithField("trace", "Env.UnlinkFtc").Error(err)
		return err
	}

	return nil
}
