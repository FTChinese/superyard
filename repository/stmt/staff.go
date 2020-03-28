package stmt

const StaffAccount = `
SELECT s.staff_id 		AS staff_id,
	IFNULL(s.email, '') AS email,
	s.user_name 		AS user_name,
	s.is_active 		AS is_active,
	s.display_name 		AS display_name,
	s.department 		AS department,
	s.group_memberships AS group_memberships`

const StaffProfile = StaffAccount + `,
	s.created_utc 				AS created_at,
	s.deactivated_utc 			AS deactivated_at,
	s.updated_utc 				AS updated_at,
	s.last_login_utc 			AS last_login_at,
	INET6_NTOA(s.last_login_ip) AS last_login_ip
FROM backyard.staff AS s`
