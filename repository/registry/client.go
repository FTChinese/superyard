package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/models/oauth"
)

const stmtCreateApp = `
INSERT INTO oauth.app_registry
SET app_name = :app_name,
	slug_name = :slug_name,
	client_id = UNHEX(:client_id),
	client_secret = UNHEX(:client_secret),
	repo_url = :repo_url,
	description = :description,
	homepage_url = :home_url,
	callback_url = :callback_url,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP(),
	owned_by = :owned_by`

// CreateApp registers a new app.
func (env Env) CreateApp(app oauth.App) error {

	_, err := env.DB.NamedExec(stmtCreateApp, app)

	if err != nil {
		logger.WithField("trace", "Env.CreateApp").Error(err)

		return err
	}

	return nil
}

const stmtListApps = stmtSelectApp + `
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

// ListApps retrieves all apps for next-api with pagination support.
func (env Env) ListApps(p gorest.Pagination) ([]oauth.App, error) {

	apps := make([]oauth.App, 0)
	err := env.DB.Select(
		&apps,
		stmtListApps,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListApps").Error(err)

		return nil, err
	}

	return apps, nil
}

const stmtApp = stmtSelectApp + `
WHERE client_id = UNHEX(?)
LIMIT 1`

// RetrieveApp retrieves an ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env Env) RetrieveApp(clientID string) (oauth.App, error) {

	var app oauth.App
	err := env.DB.Get(&app, stmtApp, clientID)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveApp").Error(err)

		return app, err
	}

	return app, nil
}

const stmtUpdateApp = `
UPDATE oauth.app_registry
SET app_name = :name,
	slug_name = :slug,
	repo_url = :repo_url,
	description = :description,
	homepage_url = :home_url,
	callback_url = :callback_url,
	updated_utc = UTC_TIMESTAMP()
WHERE client_id = UNHEX(:client_id)
	AND owned_by = :owned_by
	AND is_active = 1
LIMIT 1`

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(app oauth.App) error {

	_, err := env.DB.NamedExec(stmtUpdateApp, app)

	if err != nil {
		logger.WithField("trace", "Env.UpdateApp").Error(err)

		return err
	}

	return nil
}

const stmtRemoveApp = `
UPDATE oauth.app_registry
	SET is_active = 0
WHERE client_id = UNHEX(?)
	AND is_active = 1
LIMIT 1`

const stmtRemoveAppKeys = `
UPDATE oauth.access
	SET is_active = 0
WHERE client_id = UNHEX(:client_id)
	AND usage_type = 'app'`

// RemoveApp deactivate an ftc app.
// Ony owner can remove his apps.
// All access tokens belonging to this app should be deactivated.
func (env Env) RemoveApp(clientID string) error {
	tx, err := env.DB.Beginx()
	if err != nil {
		logger.WithField("trace", "RemoveApp").Error()
		return err
	}

	_, err = tx.Exec(stmtRemoveApp, clientID)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.RemoveApp").Error(err)
	}

	_, err = tx.Exec(stmtRemoveAppKeys, clientID)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.RemoveApp").Error(err)
	}

	if err := tx.Commit(); err != nil {
		logger.WithField("trace", "Env.RemoveApp").Error(err)
		return err
	}

	return nil
}
