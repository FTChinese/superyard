package oauth

import (
	"errors"
	"strings"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
)

// BaseAccess is the input data submitted by client.
type BaseAccess struct {
	Description null.String `json:"description" db:"description"`
	ClientID    conv.HexBin `json:"clientId" db:"client_id"`
}

func (a BaseAccess) Validate() *render.ValidationError {

	a.Description.String = strings.TrimSpace(a.Description.String)

	ve := validator.New("description").
		MaxLen(256).
		Validate(a.Description.String)

	if ve != nil {
		return ve
	}

	return nil
}

// Access is an OAuth 2.0 access Token used by an app or person to access ftc api
type Access struct {
	ID        int64       `json:"id" gorm:"primaryKey"`
	Token     conv.HexBin `json:"token" db:"token"`
	IsActive  bool        `json:"isActive" db:"is_active"`
	ExpiresIn null.Int    `json:"expiredIn" db:"expires_in"` // Output only
	Kind      KeyKind     `json:"kind" db:"usage_type"`
	BaseAccess
	CreatedBy  string      `json:"createdBy" db:"created_by"`
	CreatedAt  chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  chrono.Time `json:"updatedAt" db:"updated_at"`
	LastUsedAt chrono.Time `json:"lastUsedAt" db:"last_used_at"`
}

func (Access) TableName() string {
	return "oauth.access"
}

func (a Access) Remove() Access {
	a.IsActive = false
	a.UpdatedAt = chrono.TimeNow()

	return a
}

// NewAccess creates a new access token instance with token generated.
// Returns error if the token cannot be generated using crypto random bytes.
func NewAccess(base BaseAccess, username string) (Access, error) {
	t, err := conv.RandomHexBin(20)
	if err != nil {
		return Access{}, err
	}

	var kind = KeyKindApp

	if base.ClientID == nil {
		kind = KeyKindPersonal
	}

	return Access{
		Token:      t,
		IsActive:   true,
		ExpiresIn:  null.Int{},
		Kind:       kind,
		BaseAccess: base,
		CreatedBy:  username,
		CreatedAt:  chrono.TimeNow(),
		UpdatedAt:  chrono.TimeNow(),
	}, nil
}

type AccessList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []Access `json:"data"`
	Err  error    `json:"-"`
}

// KeyRemover specifies the where condition when removing
// an access token.
type KeyRemover struct {
	ID        int64  `db:"id"`
	CreatedBy string `json:"createdBy" db:"created_by"`
}

// KeySelector specifies the filtering condition when
// retrieving access tokens.
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
