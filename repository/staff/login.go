package staff

import "gitlab.com/ftchinese/superyard/models/employee"

func (env Env) Login(l employee.Login) (employee.Account, error) {
	var a employee.Account
	err := env.DB.Get(&a, stmtLogin, l.UserName, l.Password)

	if err != nil {
		logger.WithField("trace", "Env.Login").Error(err)

		return a, err
	}

	return a, nil
}

// UpdateLastLogin saves user login footprint after successfully authenticated.
func (env Env) UpdateLastLogin(l employee.Login, ip string) error {
	_, err := env.DB.Exec(stmtUpdateLastLogin, ip, l.UserName)

	if err != nil {
		logger.WithField("trace", "Env.UpdateLastLogin").Error(err)

		return err
	}

	return nil
}

// SavePwResetToken send a password reset token to a user's email
func (env Env) SavePwResetToken(th employee.TokenHolder) error {
	_, err := env.DB.NamedExec(stmtInsertResetToken, &th)

	if err != nil {
		logger.WithField("trace", "Env.SavePwResetToken").Error(err)

		return err
	}

	return nil
}

// VerifyResetToken finds the account associated with a password reset token
// If an account associated with a token is found, allow user to reset password.
func (env Env) LoadResetToken(token string) (employee.TokenHolder, error) {
	var th employee.TokenHolder
	err := env.DB.Get(&th, stmtSelectResetToken, token)

	if err != nil {
		logger.WithField("trace", "Env.VerifyResetToken").Error(err)

		return th, err
	}

	return th, err
}

// DeleteResetToken deletes a password reset token after it was used.
func (env Env) DeleteResetToken(token string) error {
	_, err := env.DB.Exec(stmtDeleteResetToken, token)
	if err != nil {
		logger.WithField("trace", "Env.DeleteResetToken").Error(err)

		return err
	}

	return nil
}
