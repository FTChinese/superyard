package oauth

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"strings"
)

func NewToken() (string, error) {
	token, err := gorest.RandomHex(20)

	if err != nil {
		return "", err
	}

	return token, nil
}

// APIKey is an OAuth 2.0 access Token used by an app or person to access ftc api
type Access struct {
	Key
	IsActive    bool        `json:"isActive" db:"is_active"`
	ExpiresIn   null.Int    `json:"expiredIn" db:"expires_in"`    // Output only
	FtcID       null.String `json:"ftcId" db:"ftc_id"`            // Input. Mutually exclusive with ClientID
	Description null.String `json:"description" db:"description"` // Input. Optional user input data. Max 256
	CreatedAt   chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   chrono.Time `json:"updatedAt" db:"updated_at"`
	LastUsedAt  chrono.Time `json:"lastUsedAt" db:"last_used_at"`
}

// NewAccess creates a new access token instance with token generated.
func NewAccess() (Access, error) {
	t, err := NewToken()
	if err != nil {
		return Access{}, err
	}

	return Access{
		Key: Key{
			Token: t,
		},
		IsActive:  true,
		ExpiresIn: null.Int{},
	}, nil
}

func (a *Access) Sanitize() {
	a.FtcID.String = strings.TrimSpace(a.FtcID.String)
	a.ClientID.String = strings.TrimSpace(a.ClientID.String)
	a.Description.String = strings.TrimSpace(a.Description.String)
	a.CreatedBy = strings.TrimSpace(a.CreatedBy)
}

func (a Access) Validate() *view.Reason {

	if a.Usage == KeyUsageApp {
		if a.FtcID.Valid {
			r := view.NewInvalid("access token for app should not have ftc id provided")
			r.Field = "ftc_id"
			r.Code = view.CodeInvalid
			return r
		}
		if a.ClientID.IsZero() {
			r := view.NewInvalid("access token for app must specify a client id")
			r.Field = "client_id"
			r.Code = view.CodeInvalid
			return r
		}
		return nil
	}

	if a.Usage == KeyUsagePersonal {
		if a.ClientID.Valid {
			r := view.NewInvalid("access token for personal use should not specify a client id")
			r.Field = "client_id"
			r.Code = view.CodeInvalid
			return r
		}
		return nil
	}

	r := view.NewInvalid("usage type must be one of app or personal")
	r.Field = "usage"
	r.Code = view.CodeInvalid
	return r
}

func (a Access) GetToken() string {
	if a.Token != "" {
		return a.Token
	}

	t, _ := NewToken()

	return t
}

// PersonalAccess is an api key used by a person.
type PersonalAccess struct {
	Access
	MyftEmail null.String `json:"myftEmail"` // optional user input data. The ftc account associated with this access Token.
	CreatedBy null.String `json:"createdBy"` // optional, for personal access Token.
}

// NewPersonalAccess creates a new personal access token with token generated.
func NewPersonalAccess() (PersonalAccess, error) {
	acc := PersonalAccess{}
	t, err := NewToken()
	if err != nil {
		return acc, err
	}

	acc.Token = t

	return acc, nil
}
