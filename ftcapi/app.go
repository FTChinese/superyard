package ftcapi

import (
	"fmt"
	"strings"

	"gitlab.com/ftchinese/backyard-api/util"
)

// App represents an application that needs to access ftc api
type App struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`         // required, max 255 chars.
	Slug         string `json:"slug"`         // required, unique, max 255 chars
	ClientID     string `json:"clientId"`     // required, 10 bytes
	ClientSecret string `json:"clientSecret"` // required, 32 bytes
	RepoURL      string `json:"repoUrl"`      // required, 255 chars
	Description  string `json:"description"`  // optional, 511 chars
	HomeURL      string `json:"homeUrl"`      // optional, 255 chars
	IsActive     bool   `json:"isActive"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	OwnedBy      string `json:"ownedBy"`
}

// Sanitize removes leading and trailing spaces
func (a *App) Sanitize() {
	a.Name = strings.TrimSpace(a.Name)
	a.Slug = strings.TrimSpace(a.Slug)
	a.RepoURL = strings.TrimSpace(a.RepoURL)
	a.Description = strings.TrimSpace(a.Description)
	a.HomeURL = strings.TrimSpace(a.HomeURL)
}

// Validate performas validation on incoming app.
func (a *App) Validate() util.ValidationResult {
	if r := util.ValidateIsEmpty(a.Name, "name"); r.IsInvalid {
		return r
	}
	if r := util.ValidateMaxLen(a.Name, 255, "name"); r.IsInvalid {
		return r
	}

	if r := util.ValidateIsEmpty(a.Slug, "slug"); r.IsInvalid {
		return r
	}

	if r := util.ValidateMaxLen(a.Slug, 255, "slug"); r.IsInvalid {
		return r
	}

	if r := util.ValidateIsEmpty(a.RepoURL, "repoUrl"); r.IsInvalid {
		return r
	}

	if r := util.ValidateMaxLen(a.RepoURL, 255, "repoUrl"); r.IsInvalid {
		return r
	}

	if r := util.ValidateMaxLen(a.Description, 500, "description"); r.IsInvalid {
		return r
	}

	return util.ValidateMaxLen(a.HomeURL, 120, "homeUrl")
}

// Ownership is used to transfer an app's ownership
type Ownership struct {
	SlugName string
	NewOwner string
	OldOwner string
}

// NewApp inserts a new row into oauth.app_registry table
func (env Env) NewApp(app App) error {
	clientID, err := util.RandomHex(10)
	if err != nil {
		logger.WithField("location", "Generating client id").WithError(err)
		return err
	}

	clientSecret, err := util.RandomHex(32)
	if err != nil {
		logger.WithField("location", "Generating client secret").WithError(err)

		return err
	}
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

	_, err = env.DB.Exec(query,
		app.Name,
		app.Slug,
		clientID,
		clientSecret,
		app.RepoURL,
		app.Description,
		app.HomeURL,
		app.OwnedBy,
	)

	if err != nil {
		logger.WithField("location", "Create new ftc app").Error(err)

		return err
	}

	return nil
}

// AppRoster retrieves all ftc app with pagination
func (env Env) AppRoster(page uint, rowCount uint) ([]App, error) {
	offset := (page - 1) * rowCount

	query := fmt.Sprintf(`
	%s
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`, stmtFTCApp)

	rows, err := env.DB.Query(query, rowCount, offset)

	if err != nil {
		logger.WithField("location", "Retrieve all ftc apps").Error(err)

		return nil, err
	}

	defer rows.Close()

	apps := make([]App, 0)
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
			logger.WithField("location", "Scan a row when retriving all apps").Error(err)

			continue
		}

		app.CreatedAt = util.ISO8601Formatter.FromDatetime(app.CreatedAt, nil)
		app.UpdatedAt = util.ISO8601Formatter.FromDatetime(app.UpdatedAt, nil)

		apps = append(apps, app)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "Rows iteration when retriving all ftc apps").Error(err)

		return nil, err
	}

	return apps, nil
}

// RetrieveApp retrieves a ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env Env) RetrieveApp(slugName string) (App, error) {
	query := fmt.Sprintf(`
	%s
	WHERE app_slug = ?
	LIMIT 1`, stmtFTCApp)

	var app App
	err := env.DB.QueryRow(query, slugName).Scan(
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
		logger.WithField("location", "Retrive one ftc app").Error(err)

		return app, err
	}

	app.CreatedAt = util.ISO8601Formatter.FromDatetime(app.CreatedAt, nil)
	app.UpdatedAt = util.ISO8601Formatter.FromDatetime(app.UpdatedAt, nil)

	return app, nil
}

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(slugName string, app App) error {
	query := `
	UPDATE oauth.app_registry
	SET app_name = ?,
	  	app_slug = ?,
        repo_url = ?,
        description = IFNULL(?, description),
        homepage_url = IFNULL(?, homepage_url)
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
		slugName,
		app.OwnedBy,
	)

	if err != nil {
		logger.WithField("location", "Updating a ftc app").Error(err)

		return err
	}

	return nil
}

// RemoveApp deactivate a ftc app
func (env Env) RemoveApp(slugName, owner string) error {
	query := `
	UPDATE oauth.app_registry
      	SET is_active = 0
	WHERE app_slug = ?
      	AND owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query, slugName, owner)

	if err != nil {
		logger.WithField("location", "Deactivate a ftc app").Error(err)

		return err
	}

	return nil
}

// TransferApp tranfers ownership of an app
// Before transfer, we must make sure the target owner actually exists.
func (env Env) TransferApp(o Ownership) error {
	query := `
	UPDATE oauth.app_registry
    	SET owned_by = ?
	WHERE app_slug = ?
		owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	_, err := env.DB.Exec(query,
		o.NewOwner,
		o.SlugName,
		o.OldOwner,
	)

	if err != nil {
		logger.WithField("location", "Transfer owership of a ftc app").Error(err)

		return err
	}

	return nil
}
