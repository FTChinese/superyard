package staff

import "strings"

// MyftAccount is the ftc account owned by a staff
type MyftAccount struct {
	ID    string `json:"myftId"`
	Email string `json:"myftEmail"`
	IsVIP bool   `json:"isVip"`
}

// MyftCredential contains data to login to FTC
type MyftCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Sanitize removes leading and trailing spaces
func (c *MyftCredential) Sanitize() {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)
}

// authMyft autenticate a user's myft account
// If credentials are wrong, returnw SQLNoRows error
func (env Env) authMyft(c MyftCredential) (MyftAccount, error) {
	query := `
	SELECT user_id AS myftId,
      email AS myftEmail,
      isvip AS isVip
    FROM cmstmp01.userinfo
    WHERE (email, password) = (?, MD5(?))
	LIMIT 1`

	var a MyftAccount
	err := env.DB.QueryRow(query, c.Email, c.Password).Scan(
		&a.ID,
		&a.Email,
		&a.IsVIP,
	)

	if err != nil {
		logger.WithField("location", "Verify staff myft account credentials").Error(err)

		return a, err
	}

	return a, nil
}

// saveMyft associates a staff to a myft account
// `token` column is uniquely constrained
func (env Env) saveMyft(userName string, myft MyftAccount) error {
	query := `
	INSERT INTO backyard.staff_myft
    SET staff_name = ?,
		myft_id = ?`

	_, err := env.DB.Exec(query, userName, myft.ID)

	if err != nil {
		logger.WithField("location", "Add myft account").Error(err)

		return err
	}

	return nil
}

// AddMyft authenticate a myft account and associated it with a staff account in passed.
func (env Env) AddMyft(userName string, c MyftCredential) error {
	a, err := env.authMyft(c)

	if err != nil {
		return err
	}

	err = env.saveMyft(userName, a)

	if err != nil {
		return err
	}

	return nil
}

// ListMyft lists all myft accounts owned by a staff.
func (env Env) ListMyft(userName string) ([]MyftAccount, error) {
	query := `
	SELECT u.user_id AS myftId,
      u.email AS myftEmail,
      u.isvip AS isVip
    FROM backyard.staff_myft AS s
      LEFT JOIN cmstmp01.userinfo AS u
      ON s.myft_id = u.user_id
	WHERE s.staff_name = ?`

	rows, err := env.DB.Query(query, userName)

	var accounts []MyftAccount

	if err != nil {
		logger.
			WithField("location", "Query myft accounts").
			Error(err)
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		var a MyftAccount

		err := rows.Scan(
			&a.ID,
			&a.Email,
			&a.IsVIP,
		)

		if err != nil {
			logger.
				WithField("location", "Scan myft account").
				Error(err)

			continue
		}

		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		logger.
			WithField("location", "Rows iteration").
			Error(err)

		return accounts, err
	}

	return accounts, nil
}

// DeleteMyft allows a user to delete a myft account
func (env Env) DeleteMyft(userName string, myftID string) error {
	query := `
	DELETE FROM backyard.staff_myft
	WHERE staff_name = ?
		AND myft_id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, userName, myftID)

	if err != nil {
		logger.WithField("location", "Deleting a myft account").Error(err)

		return err
	}

	return nil
}
