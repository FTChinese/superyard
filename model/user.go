package model

import (
	"database/sql"
	"fmt"
	"gitlab.com/ftchinese/backyard-api/util"
	"strings"
	"time"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/user"
)

// UserEnv handles FTC user data.
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
func (env UserEnv) loadAccount(col, val string) (user.Account, error) {
	// NOTE: in LEFT JOIN statement, the right-hand statement are null by default, regardless of their column definitions.
	query := fmt.Sprintf(`
	SELECT u.user_id AS id,
		u.wx_union_id AS uUnionId,
		u.email AS email,
		u.user_name AS userName,
	    u.is_vip AS isVip,
	    u.mobile_phone_no AS mobile,
	    u.created_utc AS createdAt,
		u.updated_utc AS updatedAt,
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
	WHERE u.%s = ?
	LIMIT 1`, col)

	var a user.Account
	var vipType int64
	var expireTime int64
	var m user.Membership

	err := env.DB.QueryRow(query, val).Scan(
		&a.UserID,
		&a.UnionID,
		&a.Email,
		&a.UserName,
		&a.IsVIP,
		&a.Mobile,
		&a.CreatedAt,
		&a.UpdatedAt,
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

// LoadAccountByEmail retrieves a user's account by email
func (env UserEnv) LoadAccountByEmail(email string) (user.Account, error) {
	return env.loadAccount(tableUser.colEmail(), email)
}

// LoadAccountByID retrieves a user's account by uuid.
func (env UserEnv) LoadAccountByID(id string) (user.Account, error) {
	return env.loadAccount(tableUser.colID(), id)
}

func (env UserEnv) ListLoginHistory(userID string, p util.Pagination) ([]user.LoginHistory, error) {
	query := `
	SELECT user_id AS userId,
		auth_method AS authMethod,
		client_type AS clientType,
		client_version AS cilentVersion,
		INET6_NTOA(user_ip) AS userIp,
		user_agent AS userAgent,
		created_utc AS createdAt
	FROM user_db.login_history
	WHERE user_id = ?
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(
		query,
		userID,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListLoginHistory").Error(err)
		return nil, err
	}

	defer rows.Close()
	var lh []user.LoginHistory

	for rows.Next() {
		var h user.LoginHistory

		err := rows.Scan(
			&h.UserID,
			&h.AuthMethod,
			&h.ClientType,
			&h.Version,
			&h.UserIP,
			&h.UserAgent,
			&h.CreatedAt,
		)
		if err != nil {
			logger.WithField("trace", "ListLoginHistory").Error(err)
			continue
		}

		lh = append(lh, h)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListLoginHistory").Error(err)
		return nil, err
	}
	return lh, nil
}

// ListOrders retrieves a user's orders that are paid successfully.
func (env UserEnv) ListOrders(userID null.String, unionID null.String, p util.Pagination) ([]user.Order, error) {
	query := fmt.Sprintf(`
	%s
	WHERE user_id IN (?, ?)
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`, stmtOrder)

	rows, err := env.DB.Query(
		query,
		userID,
		unionID,
		p.Limit,
		p.Offset())

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

// LoadWxInfo retrieves wechat user info
func (env UserEnv) LoadWxInfo(unionID string) (user.WxInfo, error) {
	query := `
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
	FROM user_db.wechat_userinfo
	WHERE union_id = ?`

	var info user.WxInfo
	var prvl string

	err := env.DB.QueryRow(query, unionID).Scan(
		&info.UnionID,
		&info.Nickname,
		&info.AvatarURL,
		&info.Gender,
		&info.Country,
		&info.Province,
		&info.City,
		&prvl,
		&info.CreatedAt,
		&info.UpdatedAt,
	)

	if err != nil {
		return info, err
	}

	info.Privileges = strings.Split(prvl, ",")

	return info, nil
}

func (env UserEnv) ListOAuthHistory(unionID string, p util.Pagination) ([]user.OAuthHistory, error) {
	query := `
	SELECT union_id AS unionId,
		open_id AS openId,
		app_id AS appId,
		client_type AS clientType,
		client_version AS clientVersion,
		user_ip AS userIp,
		user_agent AS userAgent,
		created_utc AS createdAt,
		updated_utc AS updatedAt
	FROM user_db.wechat_access
	WHERE union_id = ?
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(
		query,
		unionID,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListOAuthHistory").Error(err)
		return nil, err
	}

	defer rows.Close()
	var ah []user.OAuthHistory

	for rows.Next() {
		var h user.OAuthHistory

		err := rows.Scan(
			&h.UnionID,
			&h.OpenID,
			&h.AppID,
			&h.ClientType,
			&h.Version,
			&h.UserIP,
			&h.UserAgent,
			&h.CreatedAt,
			&h.UpdatedAt,
		)
		if err != nil {
			logger.WithField("trace", "ListOAuthHistory").Error(err)
			continue
		}

		ah = append(ah, h)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListOAuthHistory").Error(err)
		return nil, err
	}
	return ah, nil
}
