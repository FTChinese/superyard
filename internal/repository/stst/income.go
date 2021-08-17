package stst

import stats2 "github.com/FTChinese/superyard/pkg/stats"

// YearlyIncome calculates the real income of a year.
// Yearly real income means the effective range of a subscription order within the a year.
// For example, if an order spans from 2019-03-20 to 2020-03-21, only the 2019-03-20 to 2019-12-31
// contribute to this year's income.
func (env Env) YearlyIncome(y stats2.FiscalYear) (stats2.FiscalYear, error) {
	query := `
	SELECT SUM(
		DATEDIFF(
			IF(end_date > ?, ?, end_date), 
			IF(start_date < ?, ?, start_date)
		) * trade_amount / DATEDIFF(end_date, start_date)
	) AS yearlyIncome
	FROM premium.ftc_trade
	WHERE (end_date >= ?) AND (start_date <= ?)`

	err := env.dbs.Read.QueryRow(query,
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
