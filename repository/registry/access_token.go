package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/pkg/oauth"
)

// CreateToken creates an access token for app or for human,
// depending on whether ClientID if provided.
func (env Env) CreateToken(acc oauth.Access) (int64, error) {
	result, err := env.DB.NamedExec(oauth.StmtInsertToken, acc)
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

// ListAccessTokens list tokens owned by an app.
func (env Env) ListAccessTokens(clientID string, p gorest.Pagination) ([]oauth.Access, error) {
	var tokens = make([]oauth.Access, 0)

	err := env.DB.Select(
		&tokens,
		oauth.StmtListAppKeys,
		clientID,
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		logger.WithField("trace", "Env.ListAccessTokens").Error(err)
		return tokens, err
	}

	return tokens, nil
}

// ListPersonalKeys loads all key owned either by an app or by a user.
func (env Env) ListPersonalKeys(owner string, p gorest.Pagination) ([]oauth.Access, error) {
	var keys = make([]oauth.Access, 0)

	err := env.DB.Select(
		&keys,
		oauth.StmtListPersonalKeys,
		owner,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListPersonalKeys").Error(err)
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

	_, err := env.DB.NamedExec(oauth.StmtRemoveToken, k)

	if err != nil {
		logger.WithField("trace", "Env.RemoveKey").Error(err)
		return err
	}

	return nil
}
