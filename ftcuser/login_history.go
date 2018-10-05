package ftcuser

import (
	"time"

	"gitlab.com/ftchinese/backyard-api/util"
)

// LoginHistory shows a user's login footprint
type LoginHistory struct {
	AuthMethod    string `json:"authMethod"`
	ClientType    string `json:"clientType"`
	ClientVersion string `json:"clientVersion"`
	UserIP        string `json:"userIp"`
	CreatedAt     string `json:"LoggedInAt"`
}

// LoginHistory shows a user's login history
func (env Env) LoginHistory(userID string) ([]LoginHistory, error) {
	query := `
	SELECT auth_method AS authMethod,
		client_type AS clientType,
		IFNULL(client_version, '') AS clientVersion,
		IFNULL(INET6_NTOA(user_ip), '') AS userIp,
		created_utc AS createdAt
	FROM user_db.login_history
	WHERE user_id = ?`

	rows, err := env.DB.Query(query, userID)
	if err != nil {
		logger.WithField("location", "Query login history")

		return nil, err
	}
	defer rows.Close()

	var lh []LoginHistory
	for rows.Next() {
		var h LoginHistory

		err := rows.Scan(
			&h.AuthMethod,
			&h.ClientType,
			&h.ClientVersion,
			&h.UserIP,
			&h.CreatedAt,
		)

		if err != nil {
			logger.WithField("location", "Scan login history")

			continue
		}

		h.CreatedAt = util.ISO8601Formatter.FromDatetime(h.CreatedAt, time.UTC)

		lh = append(lh, h)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "Login history rows iteration").Error(err)
		return nil, err
	}

	return lh, nil
}
