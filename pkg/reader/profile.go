package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// WxProfile show wx-only account
type WxProfile struct {
	UnionID   string      `json:"unionId" db:"union_id" gorm:"primaryKey"`
	Nickname  null.String `json:"nickname" gorm:"column:nickname"`
	AvatarURL null.String `json:"avatarUrl" gorm:"column:avatar_url"`
	Gender    enum.Gender `json:"gender" gorm:"column:gender"`
	Country   null.String `json:"country" gorm:"column:country"`
	Province  null.String `json:"province" gorm:"column:province"`
	City      null.String `json:"city" gorm:"column:city"`
	CreatedAt chrono.Time `json:"createdAt" gorm:"column:created_utc"`
	UpdatedAt chrono.Time `json:"updatedAt" gorm:"column:updated_utc"`
}

func (w WxProfile) TableName() string {
	return "user_db.wechat_userinfo"
}
