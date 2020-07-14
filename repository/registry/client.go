package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/pkg/oauth"
)

// CreateApp registers a new app.
func (env Env) CreateApp(app oauth.App) error {

	_, err := env.DB.NamedExec(oauth.StmtInsertApp, app)

	if err != nil {
		logger.WithField("trace", "Env.CreateApp").Error(err)

		return err
	}

	return nil
}

// ListApps retrieves all apps for next-api with pagination support.
func (env Env) ListApps(p gorest.Pagination) ([]oauth.App, error) {

	apps := make([]oauth.App, 0)
	err := env.DB.Select(
		&apps,
		oauth.StmtListApps,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListApps").Error(err)

		return nil, err
	}

	return apps, nil
}

// RetrieveApp retrieves an ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env Env) RetrieveApp(clientID string) (oauth.App, error) {

	var app oauth.App
	err := env.DB.Get(&app, oauth.StmtApp, clientID)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveApp").Error(err)

		return app, err
	}

	return app, nil
}

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(app oauth.App) error {

	_, err := env.DB.NamedExec(oauth.StmtUpdateApp, app)

	if err != nil {
		logger.WithField("trace", "Env.UpdateApp").Error(err)

		return err
	}

	return nil
}

// RemoveApp deactivate an ftc app.
// All access tokens belonging to this app should be deactivated.
func (env Env) RemoveApp(clientID string) error {
	tx, err := env.DB.Beginx()
	if err != nil {
		logger.WithField("trace", "RemoveApp").Error()
		return err
	}

	_, err = tx.Exec(oauth.StmtRemoveApp, clientID)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.RemoveApp").Error(err)
	}

	_, err = tx.Exec(oauth.StmtRemoveAppKeys, clientID)
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
