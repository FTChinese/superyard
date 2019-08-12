package employee

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// Profile contains the full data of a staff
type Profile struct {
	Account
	CreatedAt     chrono.Time `json:"createdAt" db:"created_at"`
	DeactivatedAt chrono.Time `json:"deactivatedAt" db:"deactivated_at"`
	UpdatedAt     chrono.Time `json:"updatedAt"`
	LastLoginAt   chrono.Time `json:"lastLoginAt"`
	LastLoginIP   null.String `json:"lastLoginIp"`
}
