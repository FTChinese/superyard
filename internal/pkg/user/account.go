package user

import (
	"github.com/guregu/null"
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
	CreatedAt   string      `json:"createdAt" gorm:"column:creatdate"`
	LastLoginAt string      `json:"lastLoginAt" gorm:"column:lastlogin"`
	LastLoginIP null.String `json:"lastLoginIp" gorm:"column:last_login_ip"`
}
