package auth

import (
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
	"gorm.io/gorm"
)

// VerifyPassword verifies a staff's password
// when user tries to change password.
// ID and Password fields are required.
func (env Env) VerifyPassword(id int64, currentPass string) (user.Account, error) {
	var a user.Account
	result := env.gormDBs.Read.
		Where("id = ? AND password = MD5(?)", id, currentPass).
		First(&a)

	if result.Error != nil {
		return user.Account{}, db.ConvertGormError(result.Error)
	}

	return a, nil
}

// UpdatePassword allows user to change password.
// Generate SQL:
// UPDATE cmstmp01.managers
// SET password = MD5(?)
// WHERE username = ?
// LIMIT 1
func (env Env) UpdatePassword(holder user.Credentials) error {

	result := env.gormDBs.Write.
		Model(&user.Account{}).
		Where("username = ?", holder.UserName).
		Limit(1).
		Update("password", gorm.Expr("MD5(?)", holder.Password))

	if result.Error != nil {
		return result.Error
	}

	return nil
}
