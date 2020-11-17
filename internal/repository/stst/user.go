package stst

import (
	stats2 "github.com/FTChinese/superyard/pkg/stats"
	"github.com/jmoiron/sqlx"
)

// Env for statistic.
type Env struct {
	DB *sqlx.DB
}

func NewEnv(db *sqlx.DB) Env {
	return Env{
		DB: db,
	}
}

// DailyNewUser finds out how many new Singup everyday.
// `start` and `end` are the time range to perform statistics.
// Time format are `YYYY-MM-DD`
func (env Env) DailyNewUser(period stats2.Period) ([]stats2.SignUp, error) {

	query := `
	SELECT COUNT(*) AS userCount,
      DATE(created_utc) AS recordDate
    FROM cmstmp01.userinfo
    WHERE DATE(created_utc) BETWEEN DATE(?) AND DATE(?)
    GROUP BY DATE(created_utc)
	ORDER BY DATE(created_utc) DESC`

	rows, err := env.DB.Query(query, period.Start, period.End)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var signups []stats2.SignUp

	for rows.Next() {
		var s stats2.SignUp
		err := rows.Scan(
			&s.Count,
			&s.Date,
		)

		if err != nil {
			continue
		}

		signups = append(signups, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return signups, nil
}
