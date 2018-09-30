package adminmodel

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Env wraps a database connection
type Env struct {
	DB *sql.DB
}

var adminLogger = log.WithFields(log.Fields{
	"package": "adminmodel",
})

// VIPRoster list all vip account on ftchinese.com
func (env Env) VIPRoster() ([]MyftVIP, error) {
	query := `
	SELECT user_id AS id,
		IFNULL(user_name, '') AS name,
		email AS email,
	FROM cmstmp01.userinfo
	WHERE isvip = 1`

	rows, err := env.DB.Query(query)

	var vips []MyftVIP

	if err != nil {
		adminLogger.WithField("location", "Query myft vip accounts").Error(err)

		return vips, err
	}
	defer rows.Close()

	for rows.Next() {
		var vip MyftVIP

		err := rows.Scan(
			&vip.ID,
			&vip.Email,
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
