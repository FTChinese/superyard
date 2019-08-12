package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// WxUser shows a wechat user's bare-bone data in
// search result.
type Wechat struct {
	UnionID   string      `json:"unionId"`
	Nickname  null.String `json:"nickname"`
	CreatedAt chrono.Time `json:"createdAt"`
	UpdatedAt chrono.Time `json:"updatedAt"`
}
