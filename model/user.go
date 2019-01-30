package model

import (
	"database/sql"
	"fmt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/user"
	"time"
)

const (
	colUserName = "user_name"
)

type UserModel struct {
	DB *sql.DB
}

func normalizeMemberTier(vipType int64) enum.Tier {
	switch vipType {

	case 10:
		return enum.TierStandard

	case 100:
		return enum.TierPremium

	default:
		return enum.InvalidTier
	}
}

func normalizePayementMethod(platform int64) string {
	switch platform {
	case 1, 3:
		return "alipay"

	case 2, 4:
		return "tenpay"

	case 8:
		return "redeem_code"

	default:
		return ""
	}
}

func normalizeClientType(platform int64) string {
	switch platform {
	case 3, 4:
		return "ios"

	default:
		return "web"
	}
}

// findUser returns a User instance by retieving a user's essential data.
// When user request password reset, you need to first find this user by email;
// When user changes email, you first need to find this user by user id.
func (env UserModel) findUser(col sqlCol, value string) (user.User, error) {
	query := fmt.Sprintf(`
	SELECT user_id AS id,
		wx_union_id AS unionId,
		email AS email,
		user_name AS name
	FROM cmstmp01.userinfo
	WHERE %s = ?
	LIMIT 1`, string(col))

	var u user.User
	err := env.DB.QueryRow(query, value).Scan(
		&u.UserID,
		&u.UnionID,
		&u.Email,
		&u.UserName,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// FindUserByEmail finds a user by email
func (env UserModel) FindUserByEmail(email string) (user.User, error) {
	return env.findUser(colEmail, email)
}

func (env UserModel) FindUserByName(name string) (user.User, error) {
	return env.findUser(colUserName, name)
}

// Find a user account by either email or user id.
// When user logins, you need to find user's account by email;
// When user updated account info, you need to find account by user id.
func (env UserModel) findAccount(col sqlCol, value string) (user.Account, error) {
	// NOTE: in LEFT JOIN statement, the right-hand statement are null by default, regardless of their column definitions.
	query := `
	SELECT u.user_id AS id,
		u.wx_union_id AS uUnionId,
		u.email AS email,
		u.user_name AS userName,		mobile_phone_no AS mobile,
		IFNULL(v.vip_type, 0) AS vipType,
		IFNULL(v.expire_time, 0) AS expireTime,
		v.member_tier AS memberTier,
		v.billing_cycle AS billingCyce,
		v.expire_date AS expireDate,
		w.nickname AS nickName
	FROM cmstmp01.userinfo AS u
		LEFT JOIN premium.ftc_vip AS v
		ON u.user_id = v.vip_id
		LEFT JOIN user_db.wechat_userinfo AS w
		ON u.wx_union_id = w.union_id
	WHERE u.user_id = ?
	LIMIT 1`

	var a user.Account
	var vipType int64
	var expireTime int64
	var m user.Membership

	err := env.DB.QueryRow(query, value).Scan(
		&a.UserID,
		&a.UnionID,
		&a.Email,
		&a.UserName,
		&a.Mobile,
		&vipType,
		&expireTime,
		&m.Tier,
		&m.Cycle,
		&m.ExpireDate,
		&a.Nickname,
	)

	if err != nil {
		logger.WithField("trace", "findAccount").Error(err)

		return a, err
	}

	// If the record is using old schema
	if m.Tier == enum.InvalidTier {
		m.Tier = normalizeMemberTier(vipType)
	}

	// If expire_date column is not empty, it will be in the form 2019-07-20.
	// This is a valid ISO8601 format and do not need to be further processed.
	if m.ExpireDate.IsZero() && expireTime != 0 {
		m.ExpireDate = chrono.DateFrom(time.Unix(expireTime, 0))
	}

	a.Membership = m

	return a, nil
}

func (env UserModel) FindAccountByEmail(email string) (user.Account, error) {
	return env.findAccount(colEmail, email)
}

// LoadOrders retrieves a user's orders that are paid successfully.
func (env UserModel) LoadOrders(userID null.String, unionID null.String) ([]user.Order, error) {
	query := `
	SELECT trade_no AS orderId,
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
	    INET6_NTOA(user_ip_bin) AS userIp
	FROM premium.ftc_trade
	WHERE user_id IN (?, ?)`

	rows, err := env.DB.Query(query, userID, unionID)
	if err != nil {
		logger.WithField("trace", "LoadOrders").Error(err)
		return nil, err
	}
	defer rows.Close()

	var orders = make([]user.Order, 0)
	for rows.Next() {
		var o user.Order
		err := rows.Scan(
			&o.ID,
			&o.Tier,
			&o.Cycle,
			&o.ListPrice,
			&o.NetPrice,
			&o.PaymentMethod,
			&o.CreatedAt,
			&o.ConfirmedAt,
			&o.StartDate,
			&o.EndDate,
			&o.ClientType,
			&o.ClientVersion,
			&o.UserIP,
		)
		if err != nil {
			logger.WithField("trace", "LoadOrders").Error(err)
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "LoadOrders").Error(err)
		return nil, err
	}

	return orders, nil
}