package ftcuser

import (
	"time"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Order is a user's subscription order
type Order struct {
	OrderID       string  `json:"orderId"`
	TierToBuy     string  `json:"tierToBuy"`
	Price         float64 `json:"price"`
	TotalAmount   float64 `json:"totalAmount"`
	BillingCycle  string  `json:"billingCycle"`
	PaymentMethod string  `json:"paymentMethod"`
	ClientType    string  `json:"clientType"`
	ClientVersion string  `json:"clientVersion"`
	CreatedAt     string  `json:"createdAt"`
	ConfirmedAt   string  `json:"confirmedAt"`
	UserIP        string  `json:"userIp"`
}

// OrderRoster retrieves all orders a user placed.
func (env Env) OrderRoster(userID string) ([]Order, error) {
	query := `
	SELECT trade_no AS orderId,
		trade_platform AS platform,
		trade_subs AS subscriptionId,
		trade_price AS price,
		trade_amount AS totalAmount,
		trade_time AS confirmedTime,
		user_ip,
		IFNULL(tier_to_buy, '') AS tierToBuy,
		IFNULL(billing_cycle, '') AS billingCycle,
		IFNULL(payment_method, '') AS paymentMethod,
		IFNULL(client_type, '') AS clientType,
		IFNULL(client_version, '') AS clientVersion,
		IFNULL(created_utc, '') AS createdAt,
		IFNULL(confirmed_utc, '') AS confirmedAt,
		IFNULL(INET6_NTOA(user_ip_bin), '') AS userIp
	FROM premium.ftc_trade
	WHERE user_id = ?`

	var orders []Order

	rows, err := env.DB.Query(query, userID)
	if err != nil {
		logger.WithField("location", "Query user orders")

		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		var platform int64
		var subsriptionID int64
		var confirmedTime string
		var userIP string

		err := rows.Scan(
			&o.OrderID,
			&platform,
			&subsriptionID,
			&o.Price,
			&o.TotalAmount,
			&confirmedTime,
			&userIP,
			&o.TierToBuy,
			&o.BillingCycle,
			&o.PaymentMethod,
			&o.ClientType,
			&o.ClientVersion,
			&o.CreatedAt,
			&o.ConfirmedAt,
			&o.UserIP,
		)

		if err != nil {
			logger.WithField("location", "Scan order")

			continue
		}

		// Those are used for schema migration.
		// Once the DB schema is upgraded, they can be removed.
		if o.TierToBuy == "" {
			o.TierToBuy = normalizeMemberTier(subsriptionID)
		}

		if o.PaymentMethod == "" {
			o.PaymentMethod = normalizePayementMethod(platform)
		}

		if o.ClientType == "" {
			o.ClientType = normalizeClientType(platform)
		}

		if o.UserIP == "" {
			o.UserIP = userIP
		}

		if o.CreatedAt != "" {
			o.CreatedAt = util.ISO8601Formatter.FromDatetime(o.CreatedAt, time.UTC)
		}

		if o.ConfirmedAt == "" {
			// Use trade_time column of old schema, which is UTC+08 time
			o.ConfirmedAt = util.ISO8601Formatter.FromDatetime(confirmedTime, util.TZShanghai)
		} else {
			// Use confirmed_utc column of new schema
			o.ConfirmedAt = util.ISO8601Formatter.FromDatetime(o.ConfirmedAt, time.UTC)
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "Order rows iteration").Error(err)
		return orders, err
	}

	return orders, nil
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
