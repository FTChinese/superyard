package oauth

import (
	"errors"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/validator"
	"strings"
)

// KeyUsage tells the kind of an access token
type KeyKind string

const (
	KeyKindApp      KeyKind = "app"      // Used by an app.
	KeyKindPersonal KeyKind = "personal" // Used by human.
)

// NewToken generated an access token using crypto random bytes.
func NewToken() (string, error) {
	token, err := gorest.RandomHex(20)

	if err != nil {
		return "", err
	}

	return token, nil
}

// BaseAccess is the input data submitted by client.
type BaseAccess struct {
	Description null.String `json:"description" db:"description"`
	ClientID    null.String `json:"clientId" db:"client_id"`
}

func (a BaseAccess) Validate() *render.ValidationError {

	a.ClientID.String = strings.TrimSpace(a.ClientID.String)
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
	ID        int64    `json:"id"`
	Token     string   `json:"token" db:"token"`
	IsActive  bool     `json:"isActive" db:"is_active"`
	ExpiresIn null.Int `json:"expiredIn" db:"expires_in"` // Output only
	Kind      KeyKind  `json:"kind" db:"usage_type"`
	BaseAccess
	CreatedBy  string      `json:"createdBy" db:"created_by"`
	CreatedAt  chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  chrono.Time `json:"updatedAt" db:"updated_at"`
	LastUsedAt chrono.Time `json:"lastUsedAt" db:"last_used_at"`
}

// NewAccess creates a new access token instance with token generated.
// Returns error if the token cannot be generated using crypto random bytes.
func NewAccess(base BaseAccess, username string) (Access, error) {
	t, err := NewToken()
	if err != nil {
		return Access{}, err
	}

	var kind = KeyKindApp

	if base.ClientID.IsZero() {
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
	}, nil
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
