package ftcuser

import (
	"fmt"
	"time"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Membership contains a user's membership information
type Membership struct {
	Tier         string `json:"tier"`
	BillingCycle string `json:"billingCycle"`
	Start        string `json:"startAt"`
	Expire       string `json:"expireAt"`
}

// Profile show the details of a registered ftc user
type Profile struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Gender       string     `json:"gender"`
	FamilyName   string     `json:"familyName"`
	GivenName    string     `json:"givenName"`
	MobileNumber string     `json:"mobileNumber"`
	Birthdate    string     `json:"birthdate"`
	Address      string     `json:"address"`
	CreatedAt    string     `json:"createdAt"`
	Membership   Membership `json:"membership"`
}

// Account show the essential information of a ftc user.
// Client might show a list of accounts and uses those data to query a user's profile, orders, etc.
type Account struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// `col` is the column name by which to find an account.
// It is used in SQL `where` clause
func (env Env) findAccount(col sqlCol, value string) (Account, error) {
	query := fmt.Sprintf(`
	SELECT user_id AS id
		IFNULL(user_name, '') AS name,
		email AS email
	FROM cmstmp01.userinfo
	WHERE %s = ?
	LIMIT 1`, string(col))

	var a Account

	err := env.DB.QueryRow(query, value).Scan(
		&a.ID,
		&a.Name,
		&a.Email,
	)

	if err != nil {
		logger.WithField("location", "Find a user account").Error(err)

		return a, err
	}

	return a, nil
}

// FindUserByName tries to find a user by userName
func (env Env) FindUserByName(userName string) (Account, error) {
	return env.findAccount(colUserName, userName)
}

// FindUserByEmail tries to find a user by email
func (env Env) FindUserByEmail(email string) (Account, error) {
	return env.findAccount(colEmail, email)
}

// Profile show a user's profile
func (env Env) Profile(userID string) (Profile, error) {
	query := `
	SELECT u.user_id AS id,
		IFNULL(u.user_name, '') AS name,
		email,
		CASE
			WHEN u.title = '101' THEN 'M'
			WHEN u.title = '102' THEN 'F'
			ELSE ''
		END AS gender,
		IFNULL(u.last_name, '') AS familyName,
		IFNULL(u.first_name, '') AS givenName,
		IFNULL(u.mobile_phone_no, '') AS mobileNumber,
		u.birthdate AS birthdate,
		IFNULL(u.address, '') AS address,
		u.register_time AS createdAt,
		v.vip_type AS vipType,
		v.expire_time AS expireTime,
		IFNULL(v.member_tier, '') AS memberTier,
		IFNULL(v.billing_cycle, '') AS billingCyce,
		IFNULL(v.start_utc, '') AS startAt,
		IFNULL(v.expire_utc, '') AS expireAt
	FROM cmstmp01.userinfo AS u
		LEFT JOIN premium.ftc_vip AS v
		ON u.user_id = v.vip_id
	WHERE u.user_id = ?`

	var p Profile
	var vipType int64
	var expireTime int64
	var m Membership

	err := env.DB.QueryRow(query, userID).Scan(
		&p.ID,
		&p.Name,
		&p.Email,
		&p.Gender,
		&p.FamilyName,
		&p.GivenName,
		&p.MobileNumber,
		&p.Birthdate,
		&p.Address,
		&p.CreatedAt,
		&vipType,
		&expireTime,
		&m.Tier,
		&m.BillingCycle,
		&m.Start,
		&m.Expire,
	)

	if err != nil {
		logger.WithField("location", "Retrievin user profile").Error(err)

		return p, err
	}
	// This table uses UTC+08:00 timezone.
	// Convert to ISO8601 in UTC.
	p.CreatedAt = util.ISO8601Formatter.FromDatetime(p.CreatedAt, util.TZShanghai)

	// If the record is using old schama, then
	// m.Tier == ""
	// m.BillingCycle == ""
	// m.Start == ""
	// m.Expire == ""
	// regardless of if the user is a member or not.
	//
	// If the user is not a member, then it will always be true:
	// vipType == 0
	// expireTime == 0

	if m.Tier == "" {
		m.Tier = normalizeMemberTier(vipType)
	}

	if m.Expire == "" {
		m.Expire = normalizeExpireTime(expireTime)
	}

	if m.Start == "" {
		m.Start = normalizeStartTime(expireTime)
	}

	p.Membership = m

	return p, nil
}

func normalizeMemberTier(vipType int64) string {
	switch vipType {

	case 10:
		return "standard"

	case 100:
		return "premium"

	default:
		return "free"
	}
}

// The passed in `timestamp` is the expire_time column
func normalizeExpireTime(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}

	return time.Unix(timestamp, 0).UTC().Format(time.RFC3339)
}

// We could only deduce the start time from expire time since billing cycle is not recorded in old schema.
// Assuming the billing cycle is one year.
func normalizeStartTime(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}

	// A Time instance in UTC
	return time.Unix(timestamp, 0).UTC().AddDate(1, 0, 0).Format(time.RFC3339)
}
