package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

// LoginHistory identifies how and from where the user login
type LoginHistory struct {
	UserID     string           `json:"userId"`
	AuthMethod enum.LoginMethod `json:"loginMethod"`
	ClientApp
	CreatedAt chrono.Time `json:"createdAt"`
}
