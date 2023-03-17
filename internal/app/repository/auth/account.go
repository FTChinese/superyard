package auth

import (
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
)

// AccountByID retrieves staff account by id column.
func (env Env) AccountByID(id int64) (user.Account, error) {
	var a user.Account

	result := env.gormDBs.Read.
		Select(user.StmtAccountCols).
		First(&a, id)
	if result.Error != nil {
		return user.Account{}, db.ConvertGormError(result.Error)
	}

	return a, nil
}

// AccountByEmail loads an account when a email
// is submitted to request a password reset letter.
func (env Env) AccountByEmail(email string) (user.Account, error) {
	var a user.Account
	result := env.gormDBs.Read.
		Select(user.StmtAccountCols).
		Where(user.StmtAccountByEmail, email).
		First(&a)

	if result.Error != nil {
		return user.Account{}, db.ConvertGormError(result.Error)
	}

	return a, nil
}

// SetEmail sets the email column is missing.
func (env Env) SetEmail(a user.Account) error {
	result := env.gormDBs.Write.
		Model(&a).
		Update("email", a.Email)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateDisplayName changes display name.
func (env Env) UpdateDisplayName(a user.Account) error {
	result := env.gormDBs.Write.
		Model(&a).
		Update("fullname", a.DisplayName)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// RetrieveProfile loads a staff's profile.
func (env Env) RetrieveProfile(id int64) (user.Profile, error) {
	var p user.Profile

	result := env.gormDBs.Read.
		First(&p, id)

	if result.Error != nil {
		return p, db.ConvertGormError(result.Error)
	}

	return p, nil
}
