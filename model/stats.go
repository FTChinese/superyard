package model

import (
	"database/sql"

	"gitlab.com/ftchinese/backyard-api/stats"
)

// StatsEnv get statistics data.
type StatsEnv struct {
	DB *sql.DB
}

// DailyNewUser finds out how many new Singup everyday.
// `start` and `end` are the time range to perform statistics.
// Time format are `YYYY-MM-DD`
func (env StatsEnv) DailyNewUser(period stats.Period) ([]stats.SignUp, error) {

	query := `
	SELECT COUNT(*) AS userCount,
      DATE(created_utc) AS recordDate
    FROM cmstmp01.userinfo
    WHERE DATE(created_utc) BETWEEN DATE(?) AND DATE(?)
    GROUP BY DATE(created_utc)
	ORDER BY DATE(created_utc) DESC`

	rows, err := env.DB.Query(query, period.Start, period.End)

	if err != nil {
		logger.WithField("trace", "DailyNewUser").Error(err)
		return nil, err
	}

	defer rows.Close()

	var signups []stats.SignUp

	for rows.Next() {
		var s stats.SignUp
		err := rows.Scan(
			&s.Count,
			&s.Date,
		)

		if err != nil {
			logger.WithField("trace", "DailyNewUser")
			continue
		}

		signups = append(signups, s)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "DailyNewUser").Error(err)
		return nil, err
	}

	return signups, nil
}
