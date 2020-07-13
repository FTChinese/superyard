package staff

const StmtAccountCols = `
SELECT s.staff_id 		AS staff_id,
	s.user_name 		AS user_name,
	IFNULL(s.email, '') AS email,
	s.is_active 		AS is_active,
	s.display_name 		AS display_name,
	s.department 		AS department,
	s.group_memberships AS group_memberships`

// StmtActiveAccountByID retrieves staff account by id,
// excluding inactive one.
const StmtActiveAccountByID = StmtAccountCols + `
FROM backyard.staff AS s
WHERE s.staff_id = ?
	AND s.is_active = 1
LIMIT 1`

// StmtAccountByID retrieves a staff account by id, including
// inactive one.
const StmtAccountByID = StmtAccountCols + `
FROM backyard.staff AS s
WHERE s.staff_id = ?
LIMIT 1`

// StmtActiveAccountByEmail retrieves an active account by email.
const StmtActiveAccountByEmail = StmtAccountCols + `
FROM backyard.staff AS s
WHERE s.email = ?
	AND s.is_active = 1
LIMIT 1`

// StmtAccountByName retrieves a staff account by user_name,
// including inactive one.
const StmtAccountByName = StmtAccountCols + `
FROM backyard.staff AS s
WHERE s.user_name = ?
LIMIT 1`

// ListAccounts retrieves a list of accounts.
// Restricted to admin privilege.
const ListAccounts = StmtAccountCols + `
FROM backyard.staff AS s
ORDER BY s.user_name ASC
LIMIT ? OFFSET ?`

const StmtUpdateAccount = `
UPDATE backyard.staff
SET user_name = :user_name,
	email = :email,
	display_name = :display_name,
	department = :department,
	group_memberships = :group_memberships,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = :staff_id
	AND is_active = 1
LIMIT 1`

const StmtAddID = `
UPDATE backyard.staff
SET staff_id = :staff_id
WHERE user_name = :user_name
LIMIT 1`

const StmtSetEmail = `
UPDATE backyard.staff
SET email = :email,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = :staff_id`

const StmtUpdateDisplayName = `
UPDATE backyard.staff
SET display_name = :display_name,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = :staff_id`

const StmtDeactivate = `
UPDATE backyard.staff
  SET is_active = 0,
	deactivated_utc = UTC_TIMESTAMP()
WHERE staff_id = ?
  AND is_active = 1
LIMIT 1`

// StmtDeletePersonalKey deactivates all personal access keys
// to API when deactivating an account.
const StmtDeletePersonalKey = `
UPDATE oauth.access
	SET is_active = 0
WHERE created_by = ?`

const StmtActivate = `
UPDATE backyard.staff
  SET is_active = 1,
	updated_utc = UTC_TIMESTAMP()
WHERE staff_id = ?
  AND is_active = 0
LIMIT 1`
