package auth

import (
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
)

// Login verifies user name and password combination.
func (env Env) Login(c user.Credentials) (user.Account, error) {
	var a user.Account
	result := env.gormDBs.Read.
		Where("username = ? AND password = MD5(?)", c.UserName, c.Password).
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

// LoadPwResetSession loads a password reset
// session data by token.
// `token` is a hexdeciaml encoded string.
func (env Env) LoadPwResetSession(token string) (user.PwResetSession, error) {
	var session user.PwResetSession
	result := env.gormDBs.Read.
		Where("token = UNHEX(?)", token).
		First(&session)

	if result.Error != nil {
		return user.PwResetSession{}, db.ConvertGormError(result.Error)
	}

	return session, nil
}

// DeleteResetToken deletes a password reset token after it was used.
// Generate SQL:
// UPDATE backyard.password_reset
// SET is_used = 1
// WHERE token = UNHEX(?)
// LIMIT 1`
func (env Env) DisableResetToken(token string) error {
	result := env.gormDBs.Write.
		Model(&user.PwResetSession{}).
		Where("token = UNHEX(?)", token).
		Limit(1).
		Update("is_used", true)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
