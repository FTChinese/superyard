package ftcapi

import (
	"fmt"
	"strings"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/util"
)



// NewAPIKey creates a new row in oauth.api_key table
func (env Env) NewAPIKey(key APIKey) error {

	token, err := gorest.RandomHex(20)

	if err != nil {
		logger.WithField("location", "Generating access token")
		return err
	}

	query := `
	INSERT INTO oauth.api_key
    SET access_token = UNHEX(?),
      	description = ?,
      	myft_id = NULLIF(?, ''),
		created_by = NULLIF(?, ''),
		owned_by_app = NULLIF(?, '')`

	_, err = env.DB.Exec(query,
		token,
		key.Description,
		key.MyftID,
		key.CreatedBy,
		key.OwnedByApp,
	)

	if err != nil {
		logger.WithField("location", "Create new ftc api key").Error(err)

		return err
	}

	return nil
}

// apiKeyRoster show all api keys owned by a user or an app
func (env Env) apiKeyRoster(w whereClause, value string) ([]APIKey, error) {
	query := fmt.Sprintf(`
	SELECT id AS id,
		LOWER(HEX(access_token)) AS token,
		description AS description,
		IFNULL(myft_id, '') AS myftId,
		created_utc AS createdAt,
		updated_utc AS updatedAt,
		IFNULL(last_used_utc, '') AS lastUsedAt
	FROM oauth.api_key
	WHERE %s
		AND is_active = 1
	ORDER BY created_utc DESC`, w.String())

	rows, err := env.DB.Query(query, value)

	if err != nil {
		logger.WithField("location", "Retrieve api keys owned by a user").Error(err)

		return nil, err
	}
	defer rows.Close()

	keys := make([]APIKey, 0)
	for rows.Next() {
		var key APIKey

		err := rows.Scan(
			&key.ID,
			&key.Token,
			&key.Description,
			&key.MyftID,
			&key.CreateAt,
			&key.UpdatedAt,
			&key.LastUsedAt,
		)

		if err != nil {
			logger.WithField("location", "Scan personal api key").Error(err)

			continue
		}

		key.CreateAt = util.ISO8601UTC.FromDatetime(key.CreateAt, nil)
		key.UpdatedAt = util.ISO8601UTC.FromDatetime(key.UpdatedAt, nil)
		if key.LastUsedAt != "" {
			key.LastUsedAt = util.ISO8601UTC.FromDatetime(key.LastUsedAt, nil)
		}

		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "Retrieve personal api keys iteration").Error(err)

		return nil, err
	}

	return keys, nil
}

// PersonalAPIKeys lists all personal access tokens owned by a user.
// This version no longer show individual token separately.
func (env Env) PersonalAPIKeys(userName string) ([]APIKey, error) {
	return env.apiKeyRoster(personalAccess, userName)
}

// AppAPIKeys show all access tokens owned by an app
func (env Env) AppAPIKeys(appSlug string) ([]APIKey, error) {
	return env.apiKeyRoster(appAccess, appSlug)
}

// Remove api key(s) owned by a person or an app.
// w determines personal key or app's key;
// id determined remove a specific key or all key owned by owner. 0 to remove all; other integer value specifies the key's id.
func (env Env) deleteAPIAccess(w whereClause, id int64, owner string) error {

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
		logger.WithField("location", "Remove personal api key").Error(err)

		return err
	}

	return nil
}

// RemovePersonalAccess removes one or all access token owned by a user.
// id == 0 removes all owned by userName;
// id > 0 removes only the one with this id.
// NOTE: SQL's auto increment key starts from 1.
func (env Env) RemovePersonalAccess(id int64, userName string) error {
	return env.deleteAPIAccess(personalAccess, id, userName)
}

// RemoveAppAccess removes one or all access token owned by an app.
// id == 0 removes all owned by this app;
// id > 0 removes only the one with the specified id.
func (env Env) RemoveAppAccess(id int64, appSlug string) error {
	return env.deleteAPIAccess(appAccess, id, appSlug)
}
