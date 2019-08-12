package customer

type sqlUserCol int

const (
	sqlUserColID sqlUserCol = iota
	sqlUserColEmail
	sqlUserColUnionID
)

func (c sqlUserCol) String() string {
	names := [...]string{
		"user_id",
		"email",
		"union_id",
	}

	if c < sqlUserColID || c > sqlUserColUnionID {
		return ""
	}

	return names[c]
}

type table int

const (
	tableUser table = iota
	tableStaff
)

func (t table) colName() string {
	return "user_name"
}

func (t table) colID() string {
	switch t {
	case tableUser:
		return "user_id"
	case tableStaff:
		return "id"
	default:
		return ""
	}
}

func (t table) colEmail() string {
	return "email"
}

const (
	stmtOrder = `
	SELECT trade_no AS orderId,
		user_id AS userId,
		tier_to_buy AS tier,
		billing_cycle AS cycle,
	    trade_price AS listPrice,
	    trade_amount AS netPrice,
		payment_method AS payMethod,
		created_utc AS createdAt,
	    confirmed_utc AS confirmedAt,
		start_date AS startDate,
		end_date AS endDate,
		client_type AS clientType,
	    client_version AS clientVersion,
	    INET6_NTOA(user_ip_bin) AS userIp,
		user_agent AS userAgent
	FROM premium.ftc_trade`

	stmtWxUser = `
	SELECT union_id AS unionId,
		nickname,
		avatar_url AS avatarUrl,
		gender,
		country,
		province,
		city,
		IFNULL(privilege, '') AS privilege,
	    created_utc AS createdAt,
	    updated_utc AS updatedAt
	FROM user_db.wechat_userinfo`
)
