// Package ftcapi is in charge to app registration,
// issuing personal access token so that only approved persons and apps
// could access next-api.
package ftcapi

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Env wraps a database connection
type Env struct {
	DB *sql.DB
}

var logger = log.WithFields(log.Fields{
	"package": "ftcapi",
})

const (
	stmtFTCApp = `
	SELECT id AS id,
		app_name AS name,
    	app_slug AS slug,
    	LOWER(HEX(client_id)) AS clientId,
    	LOWER(HEX(client_secret)) AS clientSecret,
    	repo_url AS repoUrl,
    	IFNULL(description, '') AS description,
    	IFNULL(homepage_url, '') AS homeUrl,
		is_active AS isActive,
		created_utc AS createdAt,
		updated_utc AS updatedAt,
    	owned_by AS ownedBy
	FROM oauth.app_registry`
)

type whereClause int

const (
	// Where clause applied to personal access tokens only
	personalAccess whereClause = 0
	// Where clause applied to access tokens owned by an app
	appAccess whereClause = 1
)

func (w whereClause) String() string {
	clauses := [...]string{
		// Only tokens created by a user
		"created_by = ? AND owned_by_app IS NULL",
		// Only tokens belong to an app
		"owned_by_app = ?",
	}

	return clauses[w]
}
