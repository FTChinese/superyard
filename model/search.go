package model

import (
	"database/sql"
	"fmt"
	"gitlab.com/ftchinese/backyard-api/user"
)

type SearchEnv struct {
	DB *sql.DB
}

// findUser returns a User instance by retrieving a user's essential data.
// When user request password reset, you need to first find this user by email;
// When user changes email, you first need to find this user by user id.
func (env SearchEnv) findUser(col, value string) (user.User, error) {
	query := fmt.Sprintf(`
	SELECT user_id AS id,
		wx_union_id AS unionId,
		email AS email,
		user_name AS name
	FROM cmstmp01.userinfo
	WHERE %s = ?
	LIMIT 1`, col)

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
func (env SearchEnv) FindUserByEmail(email string) (user.User, error) {
	return env.findUser(
		tableUser.colEmail(),
		email)
}

func (env SearchEnv) FindUserByName(name string) (user.User, error) {
	return env.findUser(
		tableUser.colName(),
		name)
}

func (env SearchEnv) FindOrder(orderID string) (user.Order, error)  {
	query := fmt.Sprintf(`
	%s
	WHERE trade_no = ?`, stmtOrder)

	var o user.Order
	err := env.DB.QueryRow(query, orderID).Scan(
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
		&o.UserIP)

	if err != nil {
		logger.WithField("trace", "FindOrder").Error(err)
		return o, err
	}

	return o, nil
}