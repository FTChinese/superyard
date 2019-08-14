package registry

import (
	"github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"gitlab.com/ftchinese/backyard-api/models/oauth"
)

// Env wraps db.
type Env struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository/access")

// CreateApp inserts a new row into oauth.app_registry table
func (env Env) CreateApp(app oauth.App) error {

	_, err := env.DB.NamedExec(stmtCreateApp, app)

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
		stmtListApps,
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
	err := env.DB.Get(&app, stmtApp, clientID)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveApp").Error(err)

		return app, err
	}

	return app, nil
}

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(app oauth.App) error {

	_, err := env.DB.NamedExec(stmtUpdateApp, app)

	if err != nil {
		logger.WithField("trace", "Env.UpdateApp").Error(err)

		return err
	}

	return nil
}

// SearchApp retrieves an app's client_id by its slug_name.
func (env Env) SearchApp(slug string) (string, error) {

	var clientID string
	err := env.DB.Get(&clientID, stmtSearchApp, slug)
	if err != nil {
		logger.WithField("trace", "Env.SearchApp").Error(err)
		return "", err
	}

	return clientID, nil
}

// RemoveApp deactivate an ftc app.
// All access tokens belonging to this app should be deactivated.
func (env Env) RemoveApp(clientID string) error {
	tx, err := env.DB.Begin()
	if err != nil {
		logger.WithField("trace", "RemoveApp").Error()
		return err
	}

	_, err = tx.Exec(stmtRemoveApp, clientID)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "Env.RemoveApp").Error(err)
	}

	_, err = tx.Exec(stmtRemoveAppKey, clientID)
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
