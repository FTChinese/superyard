package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

var StmtAccountCols = []string{"id", "username", "email", "fullname"}

const (
	StmtAuthBy         = "username = ? AND password = MD5(?)"
	StmtVerifyPass     = "id = ? AND password = MD5(?)"
	StmtAccountByEmail = "email = ?"
)

type Account struct {
	ID          int64  `json:"id"`
	UserName    string `json:"userName" gorm:"column:username"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName" gorm:"column:fullname"`
}

func (Account) TableName() string {
	return "cmstmp01.managers"
}

// NormalizedName pick a title to name user in email.
func (a Account) NormalizeName() string {
	if a.DisplayName != "" {
		return a.DisplayName
	}

	return a.UserName
}

type Profile struct {
	Account
	CreatedAt   chrono.Time `json:"createdAt" gorm:"column:creatdate"`
	LastLoginAt chrono.Time `json:"lastLoginAt" gorm:"column:lastlogin"`
	LastLoginIP null.String `json:"lastLoginIp" gorm:"column:last_login_ip"`
}

// StmtListAccounts retrieves a list of accounts.
// Restricted to admin privilege.
const StmtListAccounts = `
SELECT s.staff_id 		AS staff_id,
	s.user_name 		AS user_name,
	IFNULL(s.email, '') AS email,
	s.is_active 		AS is_active,
	s.display_name 		AS display_name,
	s.department 		AS department,
	s.group_memberships AS group_memberships
FROM backyard.staff AS s
ORDER BY s.user_name ASC
LIMIT ? OFFSET ?`

const StmtCountStaff = `
SELECT COUNT(*)
FROM backyard.staff`
