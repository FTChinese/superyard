package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/pkg/oauth"
)

const stmtInsertToken = `
INSERT INTO oauth.access
SET access_token = UNHEX(:token),
	is_active = :is_active,
	expires_in = :expires_in,
	usage_type = :usage_type,
	description = :description,
	created_by = :created_by,
	client_id = UNHEX(:client_id),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

// CreateToken creates an access token.
func (env Env) CreateToken(acc oauth.Access) (int64, error) {
	result, err := env.DB.NamedExec(stmtInsertToken, acc)
	if err != nil {
		logger.WithField("trace", "Env.CreateKey").Error(err)

		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Retrieve keys owned by an app.
const stmtAccessTokens = stmtSelectToken + `
WHERE k.is_active = 1
	AND k.client_id = UNHEX(?)
	AND k.usage_type = 'app'
ORDER BY k.created_utc DESC
LIMIT ? OFFSET ?`

// ListAccessTokens list tokens owned by an app.
func (env Env) ListAccessTokens(clientID string, p gorest.Pagination) ([]oauth.Access, error) {
	var tokens = make([]oauth.Access, 0)

	err := env.DB.Select(&tokens, stmtAccessTokens, clientID, p.Limit, p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListAccessTokens").Error(err)
		return tokens, err
	}

	return tokens, nil
}

// Retrieve a staff's personal keys.
const stmtPersonalKeys = stmtSelectToken + `
WHERE k.is_active = 1
	AND k.created_by = ?
	AND k.usage_type = 'personal'
ORDER BY k.created_utc DESC
LIMIT ? OFFSET ?`

// ListPersonalKeys loads all key owned either by an app or by a user.
func (env Env) ListPersonalKeys(owner string, p gorest.Pagination) ([]oauth.Access, error) {
	var keys = make([]oauth.Access, 0)

	err := env.DB.Select(&keys, stmtPersonalKeys, owner, p.Limit, p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListPersonalKeys").Error(err)
		return keys, err
	}

	return keys, nil
}

// Deactivate a key by whoever created it.
const stmtRemoveKey = `
UPDATE oauth.access
	SET is_active = 0
WHERE id = :id
	AND created_by = :created_by
LIMIT 1`

// RemoveKey deactivate an access token owned by a user.
// An access token could only be deactivated by its creator,
// regardless of whether it is of kind personal or app.
// The id is collected from path parameter while
// owner name is retrieved from JWT.
func (env Env) RemoveKey(k oauth.Access) error {

	_, err := env.DB.NamedExec(stmtRemoveKey, k)

	if err != nil {
		logger.WithField("trace", "Env.RemoveKey").Error(err)
		return err
	}

	return nil
}
