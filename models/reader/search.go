package reader

import "github.com/guregu/null"

type FtcInfo struct {
	ID    string `json:"id" db:"ftc_id"`
	Email string `json:"email" db:"email"`
}

type WxInfo struct {
	UnionID  string      `json:"unionId" db:"union_id"`
	Nickname null.String `json:"nickname" db:"nickname"`
}
