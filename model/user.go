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

//
type UserEnv struct {
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

// LoadAccount retrieves a user account
func (env UserEnv) LoadAccount(userID string) (user.Account, error) {
	// NOTE: in LEFT JOIN statement, the right-hand statement are null by default, regardless of their column definitions.
	query := `
	SELECT u.user_id AS id,
		u.wx_union_id AS uUnionId,
		u.email AS email,
		u.user_name AS userName,
	    u.is_vip AS isVip,
	    u.mobile_phone_no AS mobile,
	    u.created_utc AS createdAt,
	    w.nickname AS nickName,
		IFNULL(v.vip_type, 0) AS vipType,
		IFNULL(v.expire_time, 0) AS expireTime,
		v.member_tier AS memberTier,
		v.billing_cycle AS billingCyce,
		v.expire_date AS expireDate
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

	err := env.DB.QueryRow(query, userID).Scan(
		&a.UserID,
		&a.UnionID,
		&a.Email,
		&a.UserName,
		&a.IsVIP,
		&a.Mobile,
		&a.CreatedAt,
		&a.Nickname,
		&vipType,
		&expireTime,
		&m.Tier,
		&m.Cycle,
		&m.ExpireDate)

	if err != nil {
		logger.WithField("trace", "loadAccount").Error(err)

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

// ListOrders retrieves a user's orders that are paid successfully.
func (env UserEnv) ListOrders(userID null.String, unionID null.String) ([]user.Order, error) {
	query := fmt.Sprintf(`
	%s
	WHERE user_id IN (?, ?)`, stmtOrder)

	rows, err := env.DB.Query(query, userID, unionID)
	if err != nil {
		logger.WithField("trace", "ListOrders").Error(err)
		return nil, err
	}
	defer rows.Close()

	var orders = make([]user.Order, 0)
	for rows.Next() {
		var o user.Order
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.LoginMethod,
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
			&o.UserAgent)

		if err != nil {
			logger.WithField("trace", "ListOrders").Error(err)
			return nil, err
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListOrders").Error(err)
		return nil, err
	}

	return orders, nil
}
