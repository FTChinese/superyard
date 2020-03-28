package admin

import "gitlab.com/ftchinese/superyard/models/staff"

const stmtCreateAccount = `
INSERT INTO backyard.staff
  SET staff_id = :staff_id,
	user_name = :user_name,
	email = :email,
	password = UNHEX(MD5(:password)),
	display_name = :display_name,
	department = :department,
	group_memberships = :group_memberships,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

// Create creates a new staff account
func (env Env) Create(a staff.SignUp) error {
	_, err := env.DB.NamedExec(stmtCreateAccount, &a)

	if err != nil {
		logger.WithField("trace", "Env.CreateAccount").Error(err)
		return err
	}

	return nil
}
