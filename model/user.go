package model

import (
	"database/sql"
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/types/util"
	"strings"
	"time"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/types/user"
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

// LoadFTCAccount retrieves a user account
func (env UserEnv) loadAccount(col sqlUserCol, val string) (user.Account, error) {
	// NOTE: in LEFT JOIN statement, the right-hand statement are null by default, regardless of their column definitions.
	var query string
	switch col {
	case sqlUserColEmail, sqlUserColID:
		query = fmt.Sprintf(stmtFtcAccount, col.String())

	case sqlUserColUnionID:
		query = stmtWxAccount
	}

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

		&a.Nickname,
		&vipType,
		&expireTime,
		&m.Tier,
		&m.Cycle,
		&m.ExpireDate,
		&a.CreatedAt,
		&a.UpdatedAt)

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
	return env.loadAccount(sqlUserColEmail, email)
}

// LoadAccountByID retrieves a user's account by uuid.
func (env UserEnv) LoadAccountByID(id string) (user.Account, error) {
	return env.loadAccount(sqlUserColID, id)
}

// LoadAccountByWx retrieves a user' wechat account.
func (env UserEnv) LoadAccountByWx(unionID string) (user.Account, error) {
	return env.loadAccount(sqlUserColUnionID, unionID)
}

func (env UserEnv) ListLoginHistory(userID string, p gorest.Pagination) ([]user.LoginHistory, error) {
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
func (env UserEnv) ListOrders(userID null.String, unionID null.String, p gorest.Pagination) ([]user.Order, error) {
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
	query := stmtWxUser + `
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

// ListWxUser show a list of wechat user.
func (env UserEnv) ListWxUser(p util.Pagination) ([]user.WxInfo, error) {
	query := stmtWxUser + `
	ORDER BY updated_utc DESC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(
		query,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListWxUser").Error(err)
		return nil, err
	}

	defer rows.Close()
	var wxUsers []user.WxInfo

	for rows.Next() {
		var u user.WxInfo
		var prvl string

		err := rows.Scan(
			&u.UnionID,
			&u.Nickname,
			&u.AvatarURL,
			&u.Gender,
			&u.Country,
			&u.Province,
			&u.City,
			&prvl,
			&u.CreatedAt,
			&u.UpdatedAt)

		if err != nil {
			logger.WithField("trace", "ListWxUser").Error(err)
			continue
		}

		wxUsers = append(wxUsers, u)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListWxUser").Error(err)
		return nil, err
	}

	return wxUsers, nil
}

// ListOAuthHistory shows a wechat user's login history.
func (env UserEnv) ListOAuthHistory(unionID string, p gorest.Pagination) ([]user.OAuthHistory, error) {
	query := `
	SELECT union_id AS unionId,
		open_id AS openId,
		app_id AS appId,
		client_type AS clientType,
		client_version AS clientVersion,
		INET6_NTOA(user_ip) AS userIp,
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
