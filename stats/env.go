package stats

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("package", "user")

// Env wraps db connection
type Env struct {
	DB *sql.DB
}

// Signup calculates how many new users signed up every day
type Signup struct {
	Count int    `json:"count"`
	Date  string `json:"date"`
}

// DailyNewUser finds out how many new Singup everyday.
// `start` and `end` are the time range to perform statistics.
// Time format are `YYYY-MM-DD`
func (env Env) DailyNewUser(period Period) ([]Signup, error) {

	query := `
	SELECT COUNT(*) AS userCount,
      DATE(register_time) AS recordDate
    FROM cmstmp01.userinfo
    WHERE DATE(register_time) BETWEEN DATE(?) AND DATE(?)
    GROUP BY DATE(register_time)
	ORDER BY DATE(register_time) DESC`

	rows, err := env.DB.Query(query, period.Start, period.End)

	if err != nil {
		logger.WithField("trace", "DailyNewUser").Error(err)
		return nil, err
	}

	defer rows.Close()

	var signups []Signup

	for rows.Next() {
		var s Signup
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
