package admin

import (
	"gitlab.com/ftchinese/superyard/models/staff"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

const stmtSelectProfile = stmt.StaffProfile + `
WHERE s.staff_id = ?
LIMIT 1`

// RetrieveProfile loads a staff's profile.
func (env Env) StaffProfile(id string) (staff.Profile, error) {
	var p staff.Profile

	err := env.DB.Get(&p, stmtSelectProfile, id)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveProfile").Error(err)

		return p, err
	}

	return p, nil
}
