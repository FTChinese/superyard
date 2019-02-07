package staff

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// Profile contains the full data of a staff
type Profile struct {
	Account
	IsActive      bool        `json:"isActive"`
	CreatedAt     chrono.Time `json:"createdAt"`
	DeactivatedAt chrono.Time `json:"deactivatedAt"`
	UpdatedAt     chrono.Time `json:"updatedAt"`
	LastLoginAt   chrono.Time `json:"lastLoginAt"`
	LastLoginIP   null.String `json:"lastLoginIp"`
}
