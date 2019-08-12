package registry

import (
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/models/oauth"
)

// SaveAppAccess saves an access token for an app.
func (env OAuthEnv) SaveAppAccess(acc oauth.Access, clientID string) (int64, error) {
	query := `
	INSERT INTO oauth.access
    SET access_token = UNHEX(?),
    	description = ?,
		client_id = UNHEX(?),
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	result, err := env.DB.Exec(query,
		acc.GetToken(),
		acc.Description,
		clientID,
	)

	if err != nil {
		logger.WithField("trace", "SaveAppAccess").Error(err)

		return 0, err
	}

	id, _ := result.LastInsertId()

	return id, nil
}

// ListAppAccess find all access tokens owned by an app, based on the slugified name
// of the app.
func (env OAuthEnv) ListAppAccess(slug string, p gorest.Pagination) ([]oauth.Access, error) {
	query := `
	SELECT t.id AS id,
		LOWER(HEX(t.access_token)) AS token,
		t.description AS description,
		t.created_utc AS createdAt,
		t.updated_utc AS updatedAt,
		t.last_used_utc AS lastUsedAt
	FROM oauth.access AS t
		JOIN oauth.app_registry AS a
		ON t.client_id = a.client_id
	WHERE t.is_active = 1
		AND a.slug_name = ?
	ORDER BY t.created_utc DESC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(
		query,
		slug,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListAppAccess").Error(err)

		return nil, err
	}
	defer rows.Close()

	keys := make([]oauth.Access, 0)
	for rows.Next() {
		var key oauth.Access

		err := rows.Scan(
			&key.ID,
			&key.Token,
			&key.Description,
			&key.CreatedAt,
			&key.UpdatedAt,
			&key.LastUsedAt,
		)

		if err != nil {
			logger.WithField("trace", "ListAppAccess").Error(err)

			continue
		}

		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListAppAccess").Error(err)

		return nil, err
	}

	return keys, nil
}

// RemoveAppAccess an access token owned by an app.
func (env OAuthEnv) RemoveAppAccess(clientID string, id int64) error {

	query := `
	UPDATE oauth.access
      SET is_active = 0
    WHERE client_id = UNHEX(?)
	  AND id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, clientID, id)

	if err != nil {
		logger.WithField("trace", "RemoveAppAccess").Error(err)

		return err
	}

	return nil
}

// FindMyftID tries to retrieve the user id of an email for an ftc account.
// Returns a nullable string regardless of whether the row exists.
// The returned user id is used to associated a personal access token with FTC account.
func (env OAuthEnv) FindMyftID(email null.String) null.String {
	if !email.Valid {
		return null.String{}
	}

	query := `
	SELECT user_id
	FROM cmstmp01.userinfo
	WHERE email = ?
	LIMIT 1`

	var userID null.String
	err := env.DB.QueryRow(
		query,
		email,
	).Scan(
		&userID,
	)

	if err != nil {
		logger.WithField("trace", "FindMyftID").Error(err)
	}

	return userID
}

// SavePersonalToken creates an access token for a human being.
func (env OAuthEnv) SavePersonalToken(acc oauth.PersonalAccess, myftID null.String) (int64, error) {
	query := `
	INSERT INTO oauth.access
    SET access_token = UNHEX(?),
		description = ?,
		myft_id = ?,
		created_by = ?,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	result, err := env.DB.Exec(
		query,
		acc.Token,
		acc.Description,
		myftID,
		acc.CreatedBy,
	)

	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()

	return id, nil
}

// ListPersonalTokens shows all the tokens used by a human.
func (env OAuthEnv) ListPersonalTokens(staffName string, p gorest.Pagination) ([]oauth.PersonalAccess, error) {
	query := fmt.Sprintf(`
	%s
	ORDER BY a.created_utc DESC
	LIMIT ? OFFSET ?`, stmtPersonalToken)

	rows, err := env.DB.Query(
		query,
		staffName,
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		logger.WithField("trace", "ListPersonalToken").Error(err)

		return nil, err
	}
	defer rows.Close()

	tokens := make([]oauth.PersonalAccess, 0)
	for rows.Next() {
		var t oauth.PersonalAccess

		err := rows.Scan(
			&t.ID,
			&t.Token,
			&t.Description,
			&t.MyftEmail,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.LastUsedAt,
		)

		if err != nil {
			logger.WithField("trace", "ListPersonalToken").Error(err)
			continue
		}

		tokens = append(tokens, t)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListPersonalToken").Error(err)
		return nil, err
	}

	return tokens, nil
}

// RemovePersonalToken deletes an access token used by a human.
func (env OAuthEnv) RemovePersonalToken(staffName string, id int64) error {
	query := `
	UPDATE oauth.access
		SET is_active = 0
	WHERE id = ?
		AND created_by = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, id, staffName)
	if err != nil {
		logger.WithField("trace", "RemovePersonalToken").Error(err)
		return err
	}

	return nil
}
