package repository

import (
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/models/subs"
)

// SearchEnv wraps db for search operations.
type SearchEnv struct {
	DB *sqlx.DB
}

// findUser returns a User instance by retrieving a user's essential data.
// When user request password reset, you need to first find this user by email;
// When user changes email, you first need to find this user by user id.
func (env SearchEnv) findUser(col, value string) (reader.User, error) {
	query := fmt.Sprintf(`
	%s
	WHERE %s = ?
	LIMIT 1`, stmtUser, col)

	var u reader.User
	err := env.DB.QueryRow(query, value).Scan(
		&u.UserID,
		&u.UnionID,
		&u.Email,
		&u.UserName,
		&u.IsVIP)

	if err != nil {
		return u, err
	}

	return u, nil
}

// FindUserByEmail finds a user by email
func (env SearchEnv) FindUserByEmail(email string) (reader.User, error) {
	return env.findUser(
		tableUser.colEmail(),
		email)
}

// FindUserByName searches an FTC user by name
func (env SearchEnv) FindUserByName(name string) (reader.User, error) {
	return env.findUser(
		tableUser.colName(),
		name)
}

func (env SearchEnv) FindUserByID(id string) (reader.User, error) {
	return env.findUser(tableUser.colID(), id)
}

func (env SearchEnv) FindWechat(nickname string, p gorest.Pagination) ([]reader.Wechat, error) {
	query := `
	SELECT union_id AS unionId,
		nickname,
		created_utc AS createdAt,
		updated_utc AS updatedAt
	FROM user_db.wechat_userinfo
	WHERE nickname LIKE ?
	ORDER BY nickname ASC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(
		query,
		"%"+nickname+"%",
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "FindWechat").Error(err)
		return nil, err
	}

	defer rows.Close()
	var wechat []reader.Wechat

	for rows.Next() {
		var w reader.Wechat

		err := rows.Scan(
			&w.UnionID,
			&w.Nickname,
			&w.CreatedAt,
			&w.UpdatedAt)

		if err != nil {
			logger.WithField("trace", "FindWechat").Error(err)
			continue
		}

		wechat = append(wechat, w)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "FindWechat").Error(err)
		return nil, err
	}

	return wechat, nil
}

// FindOrder searches for an subscription order.
func (env SearchEnv) FindOrder(orderID string) (reader.Order, error) {
	query := fmt.Sprintf(`
	%s
	WHERE trade_no = ?`, stmtOrder)

	var o reader.Order
	err := env.DB.QueryRow(query, orderID).Scan(
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
		logger.WithField("trace", "FindOrder").Error(err)
		return o, err
	}

	return o, nil
}

func (env SearchEnv) GiftCard(serial string) (subs.GiftCard, error) {
	query := `
	SELECT card_id AS id,
		serial_number AS serialNumber,
		DATE(FROM_UNIXTIME(expire_time)) AS expireDate,
		FROM_UNIXTIME(active_time) AS redeemedAt,
		tier AS tier,
		cycle_unit AS cycleUnit,
		cycle_value AS cycleCount
	FROM premium.scratch_card
	WHERE serial_number = ?`

	var c subs.GiftCard
	err := env.DB.QueryRow(query, serial).Scan(
		&c.ID,
		&c.Serial,
		&c.ExpireDate,
		&c.RedeemedAt,
		&c.Tier,
		&c.CycleUnit,
		&c.CycleCount)

	if err != nil {
		logger.WithField("trace", "GiftCard").Error(err)
		return c, err
	}

	return c, nil
}
