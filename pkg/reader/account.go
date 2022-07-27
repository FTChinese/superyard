package reader

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
)

// Wechat contain the essential data to identify a wechat user.
type Wechat struct {
	WxNickname  null.String `json:"nickname" db:"wx_nickname"`
	WxAvatarURL null.String `json:"avatarUrl" db:"wx_avatar_url"`
}

// FtcAccount contains ftc-only reader account data.
type FtcAccount struct {
	IDs
	StripeID   null.String `json:"stripeId" db:"stripe_id"`
	Email      null.String `json:"email" db:"email"`
	UserName   null.String `json:"userName" db:"user_name"`
	Password   string      `json:"password,omitempty" db:"password"`    // used only for sandbox user..
	CreatedBy  string      `json:"createdBy,omitempty" db:"created_by"` // Used only for sandbox user
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
	VIP        bool        `json:"vip" db:"is_vip"`
}

func (a FtcAccount) IsTest() bool {
	return strings.HasSuffix(a.Email.String, ".test@ftchinese.com") || strings.HasSuffix(a.Email.String, ".sandbox@ftchinese.com")
}

// Deprecated
type FtcAccountList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []FtcAccount `json:"data"`
	Err  error        `json:"-"`
}

// NormalizedName gets an FTC account's user name,
// and falls back to name part of email if not user name is not set.
func (a FtcAccount) NormalizedName() string {
	if a.UserName.Valid {
		return a.UserName.String
	}

	if a.Email.Valid {
		return strings.Split(a.Email.String, "@")[0]
	}

	return ""
}

// JoinedAccount contains both ftc account and wechat account.
// Kind is set to ftc is email exists, otherwise wx.
type JoinedAccount struct {
	FtcAccount
	Kind   enum.AccountKind `json:"kind"`
	Wechat Wechat           `json:"wechat"`
}

func (a *JoinedAccount) SetKind() {
	if a.FtcID.Valid {
		a.Kind = enum.AccountKindFtc
		return
	}

	a.Kind = enum.AccountKindWx
}

// Account contains a complete user account, consisting of
// both ftc account and wechat account.
type Account struct {
	JoinedAccount
	Membership Membership `json:"membership"`
}

// JoinedAccountSchema is used as SQL scan target to retrieve both ftc account and wechat account in a JOIN statement.
type JoinedAccountSchema struct {
	FtcAccount
	Wechat
}

func (s JoinedAccountSchema) JoinedAccount() JoinedAccount {
	a := JoinedAccount{
		FtcAccount: s.FtcAccount,
		Wechat:     s.Wechat,
	}

	a.SetKind()

	return a
}

func (s JoinedAccountSchema) BuildAccount(m Membership) Account {
	if s.VIP {
		m.Tier = enum.TierVIP
	}

	return Account{
		JoinedAccount: s.JoinedAccount(),
		Membership:    m,
	}
}
