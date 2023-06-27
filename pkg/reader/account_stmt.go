package reader

import "fmt"

const StmtAccountByFtc = `
SELECT u.user_id AS ftc_id,
	u.wx_union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name,
	w.nickname AS wx_nickname,
	IFNULL(u.is_vip, FALSE) AS is_vip
FROM cmstmp01.userinfo AS u
	LEFT JOIN user_db.wechat_userinfo AS w 
	ON u.wx_union_id = w.union_id
WHERE u.user_id = ?
LIMIT 1
`

const StmtAccountByWx = `
SELECT u.user_id AS ftc_id,
	w.union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name,
	w.nickname AS wx_nickname,
	IFNULL(u.is_vip, FALSE) AS is_vip
FROM user_db.wechat_userinfo AS w
	LEFT JOIN cmstmp01.userinfo AS u
	ON w.union_id = u.wx_union_id
WHERE w.union_id = ?
LIMIT 1
`

func GetAccountStmt(by SearchBy) (string, error) {
	switch by {
	case SearchByEmail:
		return StmtAccountByFtc, nil

	case SearchByWxName:
		return StmtAccountByWx, nil

	default:
		return "", fmt.Errorf("not supported account criteria: %d", by)
	}
}

const StmtSearchEmail = `
SELECT user_id AS id
FROM cmstmp01.userinfo
WHERE email = ?
LIMIT 1
`

const StmtSearchWxName = `
SELECT union_id AS id
FROM user_db.wechat_userinfo
WHERE nickname = ?
LIMIT 1
`

func GetSearchStmt(by SearchBy) (string, error) {
	switch by {
	case SearchByEmail:
		return StmtSearchEmail, nil

	case SearchByWxName:
		return StmtSearchWxName, nil

	default:
		return "", fmt.Errorf("not supported search criteria: %d", by)
	}
}
