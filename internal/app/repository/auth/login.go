package auth

import (
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/FTChinese/superyard/pkg/db"
)

// Login verifies user name and password combination.
func (env Env) Login(c user.Credentials) (user.Account, error) {
	var a user.Account
	result := env.gormDBs.Read.
		Select(user.StmtAccountCols).
		Where(user.StmtAuthBy, c.UserName, c.Password).
		First(&a)

	if result.Error != nil {
		return user.Account{}, db.ConvertGormError(result.Error)
	}

	return a, nil
}

// SavePwResetSession saves the password reset token.
// FIX: convert token to varbinary
func (env Env) SavePwResetSession(session user.PwResetSession) error {
	result := env.gormDBs.Write.
		Create(&session)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (env Env) LoadPwResetSession(token string) (user.PwResetSession, error) {
	var session user.PwResetSession
	result := env.gormDBs.Read.
		Where("token = ?", conv.HexStr(token)).
		First(&session)

	if result.Error != nil {
		return user.PwResetSession{}, db.ConvertGormError(result.Error)
	}

	return session, nil
}

// DeleteResetToken deletes a password reset token after it was used.
func (env Env) DisableResetToken(token string) error {
	result := env.gormDBs.Write.
		Exec(user.StmtDisableResetToken, token)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
