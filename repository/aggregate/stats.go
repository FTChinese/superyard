package aggregate

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/models/promo"

	"gitlab.com/ftchinese/backyard-api/models/stats"
)

// StatsEnv get statistics data.
type StatsEnv struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository")

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

// YearlyIncome calculates the real income of a year.
// Yearly real income means the effective range of a subscription order within the a year.
// For example, if an order spans from 2019-03-20 to 2020-03-21, only the 2019-03-20 to 2019-12-31
// contribute to this year's income.
func (env StatsEnv) YearlyIncome(y promo.FiscalYear) (promo.FiscalYear, error) {
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
		logger.WithField("trace", "YearlyIncome").Error(err)
		return y, err
	}

	return y, nil
}
