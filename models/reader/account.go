package reader

// Account shows a user's core data.
type Account struct {
	Ftc        FtcAccount `json:"ftc"`
	Wechat     WxAccount  `json:"wechat"`
	Membership Membership `json:"membership"`
}
