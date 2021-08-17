package registry

import (
	oauth2 "github.com/FTChinese/superyard/internal/pkg/oauth"
)

// CreateToken creates an access token for app or for human,
// depending on whether ClientID if provided.
func (env Env) CreateToken(acc oauth2.Access) (int64, error) {
	result, err := env.dbs.Write.NamedExec(oauth2.StmtInsertToken, acc)
	if err != nil {

		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ListAppTokens list tokens owned by an app.
func (env Env) ListAppTokens(clientID string) ([]oauth2.Access, error) {
	var tokens = make([]oauth2.Access, 0)

	err := env.dbs.Read.Select(
		&tokens,
		oauth2.StmtListAppKeys,
		clientID,
	)

	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

// ListPersonalKeys loads all key owned either by an app or by a user.
func (env Env) ListPersonalKeys(owner string) ([]oauth2.Access, error) {
	var keys = make([]oauth2.Access, 0)

	err := env.dbs.Read.Select(
		&keys,
		oauth2.StmtListPersonalKeys,
		owner)

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
func (env Env) RemoveKey(k oauth2.Access) error {

	_, err := env.dbs.Read.NamedExec(oauth2.StmtRemoveToken, k)

	if err != nil {
		return err
	}

	return nil
}
