package model

import (
	"database/sql"
	"fmt"
	"gitlab.com/ftchinese/backyard-api/oauth"
	"gitlab.com/ftchinese/backyard-api/util"
)

const (
	stmtFTCApp = `
	SELECT id AS id,
		app_name AS appName,
    	slug_name AS slugName,
    	LOWER(HEX(client_id)) AS clientId,
    	LOWER(HEX(client_secret)) AS clientSecret,
    	repo_url AS repoUrl,
    	description AS description,
    	homepage_url AS homeUrl,
		is_active AS isActive,
		created_utc AS createdAt,
		updated_utc AS updatedAt,
    	owned_by AS ownedBy
	FROM oauth.app_registry`

	stmtPersonalToken = `
	SELECT a.id AS id,
		LOWER(HEX(a.access_token)) AS token,
	    a.description AS description,
	    u.email AS ftcEmail,
	    a.created_by AS createdBy,
		a.created_utc AS createdAt,
		a.updated_utc AS updatedAt,
		a.last_used_utc AS lastUsedAt
	FROM oauth.access AS a
		LEFT JOIN cmstmp01.userinfo AS u
		ON a.myft_id = u.user_id
	WHERE a.is_active = 1
		AND a.created_by = ?
		AND a.client_id IS NULL`
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
	WHERE app_slug = ?
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
	  	app_slug = ?,
        repo_url = ?,
        description = ?,
        homepage_url = ?
	WHERE app_slug = ?
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
	WHERE client_id = ?
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