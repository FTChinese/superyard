package model

import (
	"fmt"
	"gitlab.com/ftchinese/backyard-api/types/staff"
	"gitlab.com/ftchinese/backyard-api/types/user"
)

// authMyft autenticate a user's myft account
// If credentials are wrong, returnw SQLNoRows error
func (env StaffEnv) authMyft(l user.Login) (user.User, error) {

	query := fmt.Sprintf(`
	%s
    WHERE (email, password) = (?, MD5(?))
	LIMIT 1`, stmtUser)

	var u user.User
	err := env.DB.QueryRow(query, l.Email, l.Password).Scan(
		&u.UserID,
		&u.UnionID,
		&u.Email,
		&u.UserName,
		&u.IsVIP)

	if err != nil {
		logger.WithField("trace", "authMyft").Error(err)
		return u, err
	}

	return u, nil
}

// saveMyft associates a staff to a myft account
// `token` column is uniquely constrained
func (env StaffEnv) saveMyft(my staff.Myft) error {
	query := `
	INSERT INTO backyard.staff_myft
    SET staff_name = ?,
		myft_id = ?,
		created_utc = UTC_TIMESTAMP()
	ON DUPLICATE KEY UPDATE staff_name = ?`

	_, err := env.DB.Exec(query,
		my.StaffName,
		my.MyftID,
		my.StaffName)

	if err != nil {
		logger.WithField("trace", "saveMyft").Error(err)
		return err
	}

	return nil
}

// AddMyft authenticate a myft account and associated it with a staff account in passed.
func (env StaffEnv) AddMyft(staffName string, l user.Login) error {
	// Verify MyftCredential is valid.
	u, err := env.authMyft(l)

	if err != nil {
		return err
	}

	err = env.saveMyft(staff.Myft{
		StaffName: staffName,
		MyftID:    u.UserID,
	})

	if err != nil {
		return err
	}

	return nil
}

// ListMyft lists all myft accounts owned by a staff.
func (env StaffEnv) ListMyft(staffName string) ([]user.User, error) {
	query := `
	SELECT u.user_id AS id,
		u.wx_union_id AS unionId,
		u.email AS email,
		u.user_name AS userName,
	    u.is_vip AS isVip
    FROM backyard.staff_myft AS s
      INNER JOIN cmstmp01.userinfo AS u
      ON s.myft_id = u.user_id
	WHERE s.staff_name = ?`

	rows, err := env.DB.Query(query, staffName)

	if err != nil {
		logger.WithField("trace", "ListMyft").Error(err)
		return nil, err
	}
	defer rows.Close()

	accounts := make([]user.User, 0)
	for rows.Next() {
		var u user.User

		err := rows.Scan(
			&u.UserID,
			&u.UnionID,
			&u.Email,
			&u.UserName,
			&u.IsVIP)

		if err != nil {
			logger.WithField("trace", "ListMyft").Error(err)
			continue
		}

		accounts = append(accounts, u)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListMyft").Error(err)
		return accounts, err
	}

	return accounts, nil
}

// DeleteMyft allows a user to delete a myft account
func (env StaffEnv) DeleteMyft(staffName, myftID string) error {
	query := `
	DELETE FROM backyard.staff_myft
	WHERE staff_name = ?
		AND myft_id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, staffName, myftID)

	if err != nil {
		logger.WithField("location", "Deleting a myft account").Error(err)

		return err
	}

	return nil
}
