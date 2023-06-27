package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// WxProfile show wx-only account
type WxProfile struct {
	UnionID   string      `json:"unionId" db:"union_id" gorm:"primaryKey"`
	Nickname  null.String `json:"nickname" db:"nickname"`
	AvatarURL null.String `json:"avatarUrl" db:"avatar_url"`
	Gender    enum.Gender `json:"gender" db:"gender"`
	Country   null.String `json:"country" db:"country"`
	Province  null.String `json:"province" db:"province"`
	City      null.String `json:"city" db:"city"`
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt chrono.Time `json:"updatedAt" db:"updated_at"`
}
