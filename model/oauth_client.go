package model

import (
	"database/sql"
	"fmt"
	"gitlab.com/ftchinese/backyard-api/oauth"
	"gitlab.com/ftchinese/backyard-api/util"
)

type OAuthEnv struct {
	DB *sql.DB
}
// NewApp inserts a new row into oauth.app_registry table
func (env OAuthEnv) SaveApp(app oauth.App) error {

	query := `
	INSERT INTO oauth.app_registry
	SET app_name = ?,
		slug_name = ?,
        client_id = UNHEX(?),
        client_secret = UNHEX(?),
        repo_url = ?,
        description = ?,
        homepage_url = ?,
		owned_by = ?`

	_, err := env.DB.Exec(query,
		app.Name,
		app.Slug,
		app.ClientID,
		app.ClientSecret,
		app.RepoURL,
		app.Description,
		app.HomeURL,
		app.OwnedBy,
	)

	if err != nil {
		logger.WithField("trace", "SaveApp").Error(err)

		return err
	}

	return nil
}

// ListApps retrieves all apps for next-api with pagination support.
func (env OAuthEnv) ListApps(p util.Pagination) ([]oauth.App, error) {

	query := fmt.Sprintf(`
	%s
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`, stmtFTCApp)

	rows, err := env.DB.Query(
		query,
		p.RowCount,
		p.Offset())

	if err != nil {
		logger.WithField("location", "Retrieve all ftc apps").Error(err)

		return nil, err
	}

	defer rows.Close()

	apps := make([]oauth.App, 0)
	for rows.Next() {
		var app oauth.App

		err := rows.Scan(
			&app.ID,
			&app.Name,
			&app.Slug,
			&app.ClientID,
			&app.ClientSecret,
			&app.RepoURL,
			&app.Description,
			&app.HomeURL,
			&app.IsActive,
			&app.CreatedAt,
			&app.UpdatedAt,
			&app.OwnedBy,
		)

		if err != nil {
			logger.WithField("trace", "ListApps").Error(err)

			continue
		}

		apps = append(apps, app)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListApps").Error(err)

		return nil, err
	}

	return apps, nil
}

// LoadApp retrieves an ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env OAuthEnv) LoadApp(slug string) (oauth.App, error) {
	query := fmt.Sprintf(`
	%s
	WHERE slug_name = ?
	LIMIT 1`, stmtFTCApp)

	var app oauth.App
	err := env.DB.QueryRow(query, slug).Scan(
		&app.ID,
		&app.Name,
		&app.Slug,
		&app.ClientID,
		&app.ClientSecret,
		&app.RepoURL,
		&app.Description,
		&app.HomeURL,
		&app.IsActive,
		&app.CreatedAt,
		&app.UpdatedAt,
		&app.OwnedBy,
	)

	if err != nil {
		logger.WithField("trace", "LoadApp").Error(err)

		return app, err
	}

	return app, nil
}

// UpdateApp allows user to update a ftc app.
func (env OAuthEnv) UpdateApp(slug string, app oauth.App) error {
	query := `
	UPDATE oauth.app_registry
	SET app_name = ?,
	  	slug_name = ?,
        repo_url = ?,
        description = ?,
        homepage_url = ?
	WHERE slug_name = ?
		AND owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query,
		app.Name,
		app.Slug,
		app.RepoURL,
		app.Description,
		app.HomeURL,
		slug,
		app.OwnedBy,
	)

	if err != nil {
		logger.WithField("trace", "UpdateApp").Error(err)

		return err
	}

	return nil
}

// FindClientID retrieves an app's client_id by its slug_name.
func (env OAuthEnv) FindClientID(slug string) (string, error) {
	query := `
	SELECT LOWER(HEX(client_id)) AS clientId
	FROM oauth.app_registry
	WHERE slug_name = ?`

	var clientID string
	err := env.DB.QueryRow(query, slug).Scan(&clientID)
	if err != nil {
		logger.WithField("trace", "FindClientID").Error(err)
		return "", err
	}

	return clientID, nil
}

// RemoveApp deactivate an ftc app.
// All access tokens belonging to this app should be deactivated.
func (env OAuthEnv) RemoveApp(clientID, owner string) error {
	tx, err := env.DB.Begin()
	if err != nil {
		logger.WithField("trace", "RemoveApp").Error()
		return err
	}

	query := `
	UPDATE oauth.app_registry
      	SET is_active = 0
	WHERE client_id = UNHEX(?)
      	AND owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	_, err = tx.Exec(query, clientID, owner)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "RemoveApp").Error(err)
	}

	query = `
	UPDATE oauth.access
    	SET is_active = 0
    WHERE client_id = ?`

	_, err = tx.Exec(query, clientID)
	if err != nil {
		_ = tx.Rollback()
		logger.WithField("trace", "RemoveApp").Error(err)
	}

	if err := tx.Commit(); err != nil {
		logger.WithField("trace", "RemoveApp").Error(err)
		return err
	}

	return nil
}

// TransferApp transfers ownership of an app.
func (env OAuthEnv) TransferApp(o oauth.Ownership) error {
	query := `
	UPDATE oauth.app_registry
    	SET owned_by = ?
	WHERE slug_name = ?
		AND owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query,
		o.NewOwner,
		o.SlugName,
		o.OldOwner,
	)

	if err != nil {
		logger.WithField("trace", "TransferApp").Error(err)

		return err
	}

	return nil
}