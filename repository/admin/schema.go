package admin

const (
	stmtInsertEmployee = `
	INSERT INTO backyard.staff
      SET user_name = :user_name,
        email = :email,
        password = UNHEX(MD5(:password)),
        display_name = :display_name,
        department = :department,
		group_memberships = :group_memberships,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	stmtSelectVIP = `
	SELECT user_id
		email,
		is_vip
	FROM cmstmp01.userinfo
	WHERE is_vip = 1
	LIMIT ? OFFSET ?`

	stmtUpdateVIP = `
	UPDATE cmstmp01.userinfo
      SET is_vip = ?
    WHERE user_id = ?
	LIMIT 1`
)
