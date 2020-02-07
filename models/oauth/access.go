package oauth

import (
	"errors"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"strings"
)

type KeyUsage string

const (
	KeyUsageApp      KeyUsage = "app"
	KeyUsagePersonal KeyUsage = "personal"
)

func NewToken() (string, error) {
	token, err := gorest.RandomHex(20)

	if err != nil {
		return "", err
	}

	return token, nil
}

type InputKey struct {
	Description null.String `json:"description"`
	CreatedBy   string      `json:"createdBy"`
	ClientID    null.String `json:"clientId"`
}

func (i InputKey) Usage() KeyUsage {
	if i.ClientID.IsZero() {
		return KeyUsagePersonal
	}

	return KeyUsageApp
}

type KeyRemover struct {
	ID        int64  `db:"id"`
	CreatedBy string `json:"createdBy" db:"created_by"`
}

type KeySelector struct {
	ClientID  null.String `query:"client_id"`
	StaffName null.String `query:"staff_name"`
}

func (s KeySelector) Validate() error {
	if s.ClientID.IsZero() && s.StaffName.IsZero() {
		return errors.New("filter criteria must be specified as either client_id or staff_name")
	}

	if s.ClientID.Valid && s.StaffName.Valid {
		return errors.New("filter by both client_id and staff_name are not supported")
	}

	return nil
}

// APIKey is an OAuth 2.0 access Token used by an app or person to access ftc api
type Access struct {
	ID          int64       `json:"id"`
	Token       string      `json:"token" db:"token"`
	IsActive    bool        `json:"isActive" db:"is_active"`
	ExpiresIn   null.Int    `json:"expiredIn" db:"expires_in"` // Output only
	Usage       KeyUsage    `json:"usage" db:"usage_type"`
	Description null.String `json:"description" db:"description"` // Input. Optional user input data. Max 256
	CreatedBy   string      `json:"created_by" db:"created_by"`
	ClientID    null.String `json:"clientId" db:"client_id"`
	CreatedAt   chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   chrono.Time `json:"updatedAt" db:"updated_at"`
	LastUsedAt  chrono.Time `json:"lastUsedAt" db:"last_used_at"`
}

// NewAccess creates a new access token instance with token generated.
func NewAccess(input InputKey) (Access, error) {
	t, err := NewToken()
	if err != nil {
		return Access{}, err
	}

	return Access{
		Token:       t,
		IsActive:    true,
		ExpiresIn:   null.Int{},
		Usage:       input.Usage(),
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
		ClientID:    input.ClientID,
		CreatedAt:   chrono.TimeNow(),
	}, nil
}

func (a *Access) Sanitize() {
	a.ClientID.String = strings.TrimSpace(a.ClientID.String)
	a.Description.String = strings.TrimSpace(a.Description.String)
	a.CreatedBy = strings.TrimSpace(a.CreatedBy)
}

func (a Access) Validate() *view.Reason {

	if a.Usage == KeyUsageApp {
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
