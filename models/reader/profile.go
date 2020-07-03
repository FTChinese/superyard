package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type Address struct {
	Country  null.String `json:"country" db:"country"`
	Province null.String `json:"province" db:"province"`
	City     null.String `json:"city" db:"city"`
	District null.String `json:"district" db:"district"`
	Street   null.String `json:"street" db:"street"`
	Postcode null.String `json:"postcode" db:"postcode"`
}

type FtcBaseProfile struct {
	Mobile    null.String `json:"mobile" db:"mobile"`
	Gender    enum.Gender `json:"gender" db:"gender"`
	LastName  null.String `json:"lastName" db:"last_name"`
	FirstName null.String `json:"firstName" db:"first_name"`
	Birthday  chrono.Date `json:"birthday" db:"birthday"`

	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt chrono.Time `json:"updatedAt" db:"updated_at"`
}

// FtcProfile show ftc-only account.
type FtcProfile struct {
	FtcAccount
	FtcBaseProfile
	Address Address `json:"address"`
}

type FtcProfileSchema struct {
	FtcAccount
	FtcBaseProfile
	Address
}

func (s FtcProfileSchema) Build() FtcProfile {
	return FtcProfile{
		FtcAccount:     s.FtcAccount,
		FtcBaseProfile: s.FtcBaseProfile,
		Address:        s.Address,
	}
}

// WxProfile show wx-only account
type WxProfile struct {
	UnionID   string      `json:"unionId" db:"union_id"`
	Nickname  null.String `json:"nickname" db:"nickname"`
	AvatarURL null.String `json:"avatarUrl" db:"avatar_url"`
	Gender    enum.Gender `json:"gender" db:"gender"`
	Country   null.String `json:"country" db:"country"`
	Province  null.String `json:"province" db:"province"`
	City      null.String `json:"city" db:"city"`
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt chrono.Time `json:"updatedAt" db:"updated_at"`
}
