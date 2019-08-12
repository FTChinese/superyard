package staff

import "gitlab.com/ftchinese/backyard-api/models/employee"

// VerifyPassword checks whether an employee's credentials are correct.
func (env Env) VerifyPassword(l employee.Login) (bool, error) {
	var matched bool
	err := env.DB.Get(&matched, stmtVerifyPassword, l.UserName, l.Password)

	if err != nil {
		logger.WithField("trace", "VerifyPassword").Error(err)
		return false, err
	}

	return matched, nil
}

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

// ResetPassword allows user to reset password after clicked the password reset link in its email.
func (env Env) ResetPassword(reset employee.PasswordReset) error {
	th, err := env.LoadResetToken(reset.Token)
	if err != nil {
		return err
	}

	profile, err := env.Load(ColumnEmail, th.Email)
	if err != nil {
		return err
	}

	if err := env.changePassword(reset.Password, profile.UserName); err != nil {
		return err
	}

	if err := env.deleteResetToken(reset.Token); err != nil {
		return err
	}

	return nil
}

// DeleteResetToken deletes a password reset token after it was used.
func (env Env) deleteResetToken(token string) error {
	_, err := env.DB.Exec(stmtDeleteResetToken, token)
	if err != nil {
		logger.WithField("trace", "Env.deleteResetToken").Error(err)

		return err
	}

	return nil
}
