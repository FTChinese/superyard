package admin

// MyftVIP is a ftc account which granted VIP
type MyftVIP struct {
	ID       string `json:"myftId"`
	Email    string `json:"myftEmail"`
	UserName string `json:"userName"`
}

// VIPRoster list all vip account on ftchinese.com
func (env Env) VIPRoster() ([]MyftVIP, error) {
	query := `
	SELECT user_id AS id,
		email AS email,
		IFNULL(user_name, '') AS name
	FROM cmstmp01.userinfo
	WHERE isvip = 1`

	rows, err := env.DB.Query(query)

	if err != nil {
		adminLogger.WithField("location", "Query myft vip accounts").Error(err)

		return nil, err
	}
	defer rows.Close()

	vips := make([]MyftVIP, 0)
	for rows.Next() {
		var vip MyftVIP

		err := rows.Scan(
			&vip.ID,
			&vip.Email,
			&vip.UserName,
		)

		if err != nil {
			adminLogger.WithField("location", "Scan myft vip account").Error(err)

			continue
		}

		vips = append(vips, vip)
	}

	if err := rows.Err(); err != nil {
		adminLogger.WithField("location", "rows iteration").Error(err)

		return vips, err
	}

	return vips, nil
}

func (env Env) updateVIP(myftID string, isVIP bool) error {
	query := `
	UPDATE cmstmp01.userinfo
      SET isvip = ?
    WHERE user_id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, isVIP, myftID)

	if err != nil {
		adminLogger.WithField("location", "Grant vip to a ftc account")

		return err
	}

	return nil
}

// GrantVIP set a ftc account as vip
func (env Env) GrantVIP(myftID string) error {
	return env.updateVIP(myftID, true)
}

// RevokeVIP removes vip status from a ftc account
func (env Env) RevokeVIP(myftID string) error {
	return env.updateVIP(myftID, false)
}
