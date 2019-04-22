package model

import (
	"database/sql"
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/user"
)

// SearchEnv wraps db for search operations.
type SearchEnv struct {
	DB *sql.DB
}

// findUser returns a User instance by retrieving a user's essential data.
// When user request password reset, you need to first find this user by email;
// When user changes email, you first need to find this user by user id.
func (env SearchEnv) findUser(col, value string) (user.User, error) {
	query := fmt.Sprintf(`
	%s
	WHERE %s = ?
	LIMIT 1`, stmtUser, col)

	var u user.User
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
func (env SearchEnv) FindUserByEmail(email string) (user.User, error) {
	return env.findUser(
		tableUser.colEmail(),
		email)
}

// FindUserByName searches an FTC user by name
func (env SearchEnv) FindUserByName(name string) (user.User, error) {
	return env.findUser(
		tableUser.colName(),
		name)
}

func (env SearchEnv) FindUserByID(id string) (user.User, error) {
	return env.findUser(tableUser.colID(), id)
}

func (env SearchEnv) FindWechat(nickname string, p gorest.Pagination) ([]user.Wechat, error) {
	query := `
	SELECT union_id AS unionId,
		nickname,
		created_utc AS createdAt,
		updated_utc AS updatedAt
	FROM user_db.wechat_info
	WHERE nickname = ?
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(
		query,
		nickname,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "FindWechat").Error(err)
		return nil, err
	}

	defer rows.Close()
	var wechat []user.Wechat

	for rows.Next() {
		var w user.Wechat

		err := rows.Scan(
			w.UnionID,
			w.Nickname,
			w.CreatedAt,
			w.UpdatedAt)

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
func (env SearchEnv) FindOrder(orderID string) (user.Order, error) {
	query := fmt.Sprintf(`
	%s
	WHERE trade_no = ?`, stmtOrder)

	var o user.Order
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

func (env SearchEnv) FindWxUser(unionID string) {

}
