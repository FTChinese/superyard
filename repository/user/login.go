package user

import (
	"gitlab.com/ftchinese/superyard/pkg/staff"
)

// Login verifies user name and password combination.
func (env Env) Login(l staff.Credentials) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, staff.StmtLogin, l.UserName, l.Password)

	if err != nil {
		logger.WithField("trace", "Env.Login").Error(err)

		return a, err
	}

	return a, nil
}

// UpdateLastLogin saves user login footprint after successfully authenticated.
func (env Env) UpdateLastLogin(l staff.Credentials, ip string) error {
	_, err := env.DB.Exec(staff.StmtUpdateLastLogin, ip, l.UserName)

	if err != nil {
		logger.WithField("trace", "Env.UpdateLastLogin").Error(err)

		return err
	}

	return nil
}

// SavePwResetSession saves the password reset token.
func (env Env) SavePwResetSession(session staff.PwResetSession) error {
	_, err := env.DB.NamedExec(staff.StmtInsertPwResetSession, session)

	if err != nil {
		logger.WithField("trace", "Env.SavePwResetSession").Error(err)

		return err
	}

	return nil
}

func (env Env) LoadPwResetSession(token string) (staff.PwResetSession, error) {
	var session staff.PwResetSession
	err := env.DB.Get(&session, staff.StmtPwResetSession, token)

	if err != nil {
		logger.WithField("trace", "Env.LoadPwResetSession").Error(err)

		return staff.PwResetSession{}, err
	}

	return session, nil
}

// AccountByResetToken finds an account by a password reset token.
// When a user submitted token and password when trying to
// reset password, we should use the token to find out
// the account of this user before updating the password.
func (env Env) AccountByResetToken(token string) (staff.Account, error) {
	var a staff.Account
	err := env.DB.Get(&a, staff.StmtAccountByResetToken, token)

	if err != nil {
		logger.WithField("trace", "Env.AccountByResetToken").Error(err)

		return staff.Account{}, err
	}

	return a, err
}

// DeleteResetToken deletes a password reset token after it was used.
func (env Env) DisableResetToken(token string) error {
	_, err := env.DB.Exec(staff.StmtDisableResetToken, token)
	if err != nil {
		logger.WithField("trace", "Env.DeleteResetToken").Error(err)

		return err
	}

	return nil
}
