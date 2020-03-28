package admin

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

const stmtListStaff = stmt.StaffAccount + `
FROM backyard.staff AS s
ORDER BY s.user_name ASC
LIMIT ? OFFSET ?`

func (env Env) ListStaff(p util.Pagination) ([]employee.Profile, error) {
	profiles := make([]employee.Profile, 0)

	err := env.DB.Select(&profiles,
		stmtListStaff,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListStaff").Error(err)

		return profiles, err
	}

	return profiles, nil
}
