package ftcapp

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Env wraps a database connection
type Env struct {
	DB *sql.DB
}

var appLogger = log.WithFields(log.Fields{
	"package": "appmodel",
})

const (
	stmtFTCApp = `
	SELECT id AS id,
		app_name AS name
    	app_slug AS slug
    	LOWER(HEX(app.client_id)) AS clientId,
    	LOWER(HEX(app.client_secret)) AS clientSecret,
    	repo_url AS repoUrl,
    	description AS description,
    	homepage_url AS homeUrl,
		is_active AS isActive,
		created_utc AS createdAt,
		updated_utc AS updatedAt,
    	owned_by AS ownedBy
	FROM app_registry`
)

// NewApp inserts a new row into oauth.app_registry table
func (env Env) NewApp(app App) error {
	query := `
	INSERT INTO oauth.app_registry
	SET app_name = ?,
		app_slug = ?,
        client_id = UNHEX(?),
        client_secret = UNHEX(?),
        repo_url = ?,
        description = NULLIF(?, ''),
        homepage_url = NULLIF(?, ''),
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
		appLogger.WithField("location", "Create new ftc app").Error(err)

		return err
	}

	return nil
}

// RetrieveApp retrieves a ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env Env) RetrieveApp(slug string) (App, error) {
	query := fmt.Sprintf(`
	%s
	WHERE app_slug = ?
	LIMIT 1`, stmtFTCApp)

	var app App
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
		appLogger.WithField("location", "Retrive one ftc app").Error(err)

		return app, err
	}

	return app, nil
}

// AppRoster retrieves all ftc app with pagination
func (env Env) AppRoster(page int, rowCount int) ([]App, error) {
	offset := (page - 1) * rowCount

	query := fmt.Sprintf(`
	%s
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`, stmtFTCApp)

	rows, err := env.DB.Query(query, rowCount, offset)

	var apps []App

	if err != nil {
		appLogger.WithField("location", "Retrieve all ftc apps").Error(err)

		return apps, err
	}

	defer rows.Close()

	for rows.Next() {
		var app App

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
			appLogger.WithField("location", "Scan a row when retriving all apps").Error(err)

			continue
		}

		apps = append(apps, app)
	}

	if err := rows.Err(); err != nil {
		appLogger.WithField("location", "Rows iteration when retriving all ftc apps").Error(err)

		return apps, err
	}

	return apps, nil
}

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(app App) error {
	query := `
	UPDATE oauth.app_registry
	  SET app_name = ?,
	  	app_slug = ?
        repo_url = ?,
        description = IFNULL(?, description),
        homepage_url = IFNULL(?, homepage_url),
    WHERE id = ?
      AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query,
		app.Name,
		app.Slug,
		app.RepoURL,
		app.Description,
		app.HomeURL,
		app.ID,
	)

	if err != nil {
		appLogger.WithField("location", "Updating a ftc app").Error(err)

		return err
	}

	return nil
}

// TransferApp tranfers ownership of an app
func (env Env) TransferApp(o Ownership) error {
	query := `
	UPDATE oauth.app_registry
    	SET owned_by = ?
	WHERE id = ?
		owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query,
		o.NewOwner,
		o.ID,
		o.OldOwner,
	)

	if err != nil {
		appLogger.WithField("location", "Transfer owership of a ftc app").Error(err)

		return err
	}

	return nil
}

// RemoveApp deactivate a ftc app
func (env Env) RemoveApp(app App) error {
	query := `
	UPDATE oauth.app_registry
      	SET is_active = 0
	WHERE id = ?
		AND app_slug = ?
      	AND owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, app.ID, app.Slug, app.OwnedBy)

	if err != nil {
		appLogger.WithField("location", "Deactivate a ftc app").Error(err)

		return err
	}

	return nil
}

// NewAPIKey creates a new row in oauth.api_key table
func (env Env) NewAPIKey(key APIKey) error {
	query := `
	INSERT INTO oauth.api_key
    SET access_token = UNHEX(?),
      	description = ?,
      	myft_id = NULLIF(?, ''),
		created_by = NULLIF(?, ''),
		owned_by_app = NULLIF(?, '')`

	_, err := env.DB.Exec(query,
		key.Token,
		key.Description,
		key.MyftID,
		key.CreatedBy,
		key.OwnedByApp,
	)

	if err != nil {
		appLogger.WithField("location", "Create new ftc api key").Error(err)

		return err
	}

	return nil
}

// APIKeyRoster show all api keys owned by a user or an app
func (env Env) APIKeyRoster(userName string, appSlug string) ([]APIKey, error) {
	var whereClause string
	if userName != "" {
		whereClause = "created_by = ? AND owned_by_app IS NULL"
	} else if appSlug != "" {
		whereClause = "owned_by_app = ?"
	}

	query := fmt.Sprintf(`
	SELECT id AS id,
		access_token AS token,
		description AS description,
		myft_id AS myftId,
		created_utc AS createdAt,
		updated_utc AS updatedAt,
		last_used_utc AS lastUsedAt
	FROM oauth.api_key
	WHERE %s
		AND is_active = 1
	ORDER BY created_utc DESC`, whereClause)

	rows, err := env.DB.Query(query, userName)

	var keys []APIKey

	if err != nil {
		appLogger.WithField("location", "Retrieve api keys owned by a user").Error(err)

		return keys, err
	}
	defer rows.Close()

	for rows.Next() {
		var key APIKey

		err := rows.Scan(
			&key.ID,
			&key.Token,
			&key.Description,
			&key.MyftID,
			&key.CreateAt,
			&key.UpdatedAt,
			&key.LastUsed,
		)

		if err != nil {
			appLogger.WithField("location", "Scan personal api key").Error(err)

			continue
		}

		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		appLogger.WithField("location", "Retrieve personal api keys iteration").Error(err)

		return keys, err
	}

	return keys, nil
}

// RemoveAPIKey removes an access token owned by a user
func (env Env) RemoveAPIKey(key APIKey) error {
	query := `
	UPDATE oauth.api_key
      SET is_active = 0
    WHERE id = ?
	  AND created_by = ?
	  AND owned_by_app IS NULL
	LIMIT 1`

	_, err := env.DB.Exec(query, key.ID, key.CreatedBy)

	if err != nil {
		appLogger.WithField("location", "Remove personal api key").Error(err)

		return err
	}

	return nil
}

type whereClause int

const (
	userKeyOnly whereClause = 0
	appKeyOnly  whereClause = 1
)

func (w whereClause) String() string {
	clauses := [...]string{
		"created_by = ? AND owned_by_app IS NULL",
		"owned_by_app = ?",
	}

	return clauses[w]
}

// Remove api key(s) owned by a person or an app.
// w determines personal key or app's key;
// id determined remove a specific key or all key owned by owner. 0 to remove all; other integer value specifies the key's id.
func (env Env) removeAPIKey(w whereClause, id int, owner string) error {

	var whereID string

	if id > 0 {
		whereID = "AND id = ?"
	}
	query := fmt.Sprintf(`
	UPDATE oauth.api_key
      SET is_active = 0
    WHERE %s
	  %s
	LIMIT 1`, w.String(), whereID)

	var err error

	if id > 0 {
		_, err = env.DB.Exec(query, owner, id)
	} else {
		_, err = env.DB.Exec(query, owner)
	}

	if err != nil {
		appLogger.WithField("location", "Remove personal api key").Error(err)

		return err
	}

	return nil
}
