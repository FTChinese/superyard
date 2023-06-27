package registry

import (
	"github.com/FTChinese/superyard/internal/pkg/oauth"
)

// CreateToken creates an access token for app or for human,
// depending on whether ClientID if provided.
func (env Env) CreateToken(acc oauth.Access) (oauth.Access, error) {

	err := env.gormDBs.Write.Create(&acc).Error

	if err != nil {
		return oauth.Access{}, err
	}

	return acc, nil
}

func (env Env) RetrieveToken(id int64, owner string) (oauth.Access, error) {
	var token oauth.Access
	err := env.gormDBs.Read.
		Where("created_by = ?", owner).
		First(&token, id).
		Error
	if err != nil {
		return oauth.Access{}, err
	}

	return token, nil

}

// ListAppTokens list tokens owned by an app.
func (env Env) ListAppTokens(clientID string) ([]oauth.Access, error) {
	var tokens = make([]oauth.Access, 0)

	err := env.gormDBs.Read.
		Where("is_active = ? AND client_id = UNHEX(?) AND usage_type = ?", true, clientID, "app").
		Order("created_utc DESC").
		Find(&tokens).
		Error

	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

// ListPersonalKeys loads all key owned either by an app or by a user.
func (env Env) ListPersonalKeys(owner string) ([]oauth.Access, error) {
	var keys = make([]oauth.Access, 0)

	err := env.gormDBs.Read.
		Where("is_active = ? AND created_by = ? AND usage_type = ?", true, owner, "personal").
		Order("created_utc DESC").
		Find(&keys).
		Error

	if err != nil {
		return keys, err
	}

	return keys, nil
}

// RemoveKey deactivate an access token owned by a user.
// An access token could only be deactivated by its creator,
// regardless of whether it is of kind personal or app.
// The id is collected from path parameter while
// owner name is retrieved from JWT.
func (env Env) RemoveKey(k oauth.Access) error {

	err := env.gormDBs.Write.
		Limit(1).
		Where("created_by = ?", k.CreatedBy).
		Save(&k).
		Error

	if err != nil {
		return err
	}

	return nil
}
