package model

import "gitlab.com/ftchinese/backyard-api/staff"

// authMyft autenticate a user's myft account
// If credentials are wrong, returnw SQLNoRows error
func (env StaffEnv) authMyft(c staff.MyftCredential) (staff.MyftAccount, error) {
	query := `
	SELECT user_id AS myftId,
      email AS myftEmail,
      isvip AS isVip
    FROM cmstmp01.userinfo
    WHERE (email, password) = (?, MD5(?))
	LIMIT 1`

	var a staff.MyftAccount
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
func (env StaffEnv) saveMyft(userName string, myft staff.MyftAccount) error {
	query := `
	INSERT INTO backyard.staff_myft
    SET staff_name = ?,
		myft_id = ?
	ON DUPLICATE KEY UPDATE staff_name = ?`

	_, err := env.DB.Exec(query, userName, myft.ID, userName)

	if err != nil {
		logger.WithField("location", "Add myft account").Error(err)

		return err
	}

	return nil
}

// AddMyft authenticate a myft account and associated it with a staff account in passed.
func (env StaffEnv) AddMyft(userName string, c staff.MyftCredential) error {
	// Verify MyftCredential is valid.
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
func (env StaffEnv) ListMyft(userName string) ([]staff.MyftAccount, error) {
	query := `
	SELECT u.user_id AS myftId,
      u.email AS myftEmail,
      u.isvip AS isVip
    FROM backyard.staff_myft AS s
      INNER JOIN cmstmp01.userinfo AS u
      ON s.myft_id = u.user_id
	WHERE s.staff_name = ?`

	rows, err := env.DB.Query(query, userName)

	if err != nil {
		logger.
			WithField("location", "Query myft accounts").
			Error(err)
		return nil, err
	}
	defer rows.Close()

	accounts := make([]staff.MyftAccount, 0)
	for rows.Next() {
		var a staff.MyftAccount

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
func (env StaffEnv) DeleteMyft(userName string, myftID string) error {
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