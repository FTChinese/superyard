package aggregate

import (
	stats2 "github.com/FTChinese/superyard/pkg/stats"
	"github.com/jmoiron/sqlx"
)

// StatsEnv get statistics data.
type StatsEnv struct {
	DB *sqlx.DB
}

// DailyNewUser finds out how many new Singup everyday.
// `start` and `end` are the time range to perform statistics.
// Time format are `YYYY-MM-DD`
func (env StatsEnv) DailyNewUser(period stats2.Period) ([]stats2.SignUp, error) {

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

// YearlyIncome calculates the real income of a year.
// Yearly real income means the effective range of a subscription order within the a year.
// For example, if an order spans from 2019-03-20 to 2020-03-21, only the 2019-03-20 to 2019-12-31
// contribute to this year's income.
func (env StatsEnv) YearlyIncome(y stats2.FiscalYear) (stats2.FiscalYear, error) {
	query := `
	SELECT SUM(
		DATEDIFF(
			IF(end_date > ?, ?, end_date), 
			IF(start_date < ?, ?, start_date)
		) * trade_amount / DATEDIFF(end_date, start_date)
	) AS yearlyIncome
	FROM premium.ftc_trade
	WHERE (end_date >= ?) AND (start_date <= ?)`

	err := env.DB.QueryRow(query,
		y.LastDate,
		y.LastDate,
		y.StartDate,
		y.StartDate,
		y.StartDate,
		y.LastDate).Scan(&y.Income)

	if err != nil {
		return y, err
	}

	return y, nil
}
