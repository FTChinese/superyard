package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"gitlab.com/ftchinese/backyard-api/types/util"
)

// LoginHistory identifies how and from where the user login
type LoginHistory struct {
	UserID     string           `json:"userId"`
	AuthMethod enum.LoginMethod `json:"loginMethod"`
	util.ClientApp
	CreatedAt chrono.Time `json:"createdAt"`
}
