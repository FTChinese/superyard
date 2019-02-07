package model

import (
	"fmt"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/oauth"
	"gitlab.com/ftchinese/backyard-api/util"
)

// SaveAppAccess saves an access token for an app.
func (env OAuthEnv) SaveAppAccess(token, clientID string) error {
	query := `
	INSERT INTO oauth.access
    SET access_token = UNHEX(?),
		client_id = ?`

	_, err := env.DB.Exec(query,
		token,
		clientID,
	)

	if err != nil {
		logger.WithField("trace", "SaveAppAccess").Error(err)

		return err
	}

	return nil
}

// ListAppAccess find all access tokens owned by an app, based on the slugified name
// of the app.
func (env OAuthEnv) ListAppAccess(slug string, p util.Pagination) ([]oauth.Access, error) {
	query := `
	SELECT t.id AS id,
		LOWER(HEX(t.access_token)) AS token,
		t.created_utc AS createdAt,
		t.updated_utc AS updatedAt,
		t.last_used_utc AS lastUsedAt
	FROM oauth.access AS t
		JOIN oauth.app_registry AS a
		ON t.client_id = a.client_id
	WHERE t.is_active = 1
		AND a.slug_name = ?
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(
		query,
		slug,
		p.RowCount,
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

// Remove an access token owned by an app.
func (env OAuthEnv) RemoveAppAccess(clientID string, id int64) error {

	query := `
	UPDATE oauth.access
      SET is_active = 0
    WHERE client_id = ?
	  AND id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, clientID)

	if err != nil {
		logger.WithField("trace", "RemoveAppAccess").Error(err)

		return err
	}

	return nil
}

// FindMyftID tries to retrieve the user id of an email for an ftc account.
// Returns a nullable string regardless of whether the row exists.
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
func (env OAuthEnv) SavePersonalToken(acc oauth.PersonalAccess, myftID null.String) error {
	query := `
	INSERT INTO oauth.access
    SET access_token = UNHEX(?),
		description = ?,
		myft_id = ?,
		created_by = ?`

	_, err := env.DB.Exec(
		query,
		acc.Token,
		acc.Description,
		myftID,
		acc.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}

// ListPersonalTokens shows all the tokens used by a human.
func (env OAuthEnv) ListPersonalTokens(userName string, p util.Pagination) ([]oauth.PersonalAccess, error) {
	query := fmt.Sprintf(`
	%s
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`, stmtPersonalToken)

	rows, err := env.DB.Query(
		query,
		userName,
		p.RowCount,
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
func (env OAuthEnv) RemovePersonalToken(userName string, id int64) error  {
	query := `
	UPDATE oauth.access
		SET is_active = 0
	WHERE id = ?
		AND created_by = ?
	LIMIT 1`
	
	_, err := env.DB.Exec(query, id, userName)
	if err != nil {
		logger.WithField("trace", "RemovePersonalToken").Error(err)
		return err
	}

	return nil
}