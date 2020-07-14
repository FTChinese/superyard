package staff

const StmtCreateAccount = `
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
