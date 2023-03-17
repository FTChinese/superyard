package auth

import (
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
)

// VerifyPassword verifies a staff's password
// when user tries to change password.
// ID and Password fields are required.
func (env Env) VerifyPassword(id int64, pass user.ParamsPasswords) (user.Account, error) {
	var a user.Account
	result := env.gormDBs.Read.
		Select(user.StmtAccountCols).
		Where(user.StmtVerifyPass, id, pass.OldPassword).
		First(&a)

	if result.Error != nil {
		return user.Account{}, db.ConvertGormError(result.Error)
	}

	return a, nil
}

// UpdatePassword allows user to change password.
// It also updates the legacy table, which does
// not have a staff_id column. So we use user_name
// to update the legacy table.
// Therefore, to update password, we should know
// user'd id and user name.
func (env Env) UpdatePassword(holder user.Credentials) error {

	result := env.gormDBs.Write.
		Exec(user.StmtUpdatePassword, holder.Password, holder.UserName)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
