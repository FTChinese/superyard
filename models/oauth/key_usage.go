package oauth

import "github.com/guregu/null"

type KeyUsage string

const (
	KeyUsageApp      KeyUsage = "app"
	KeyUsagePersonal KeyUsage = "personal"
)

// KeyOwner contains data about the key owner.
// Owner might be an app or a human being.
type KeyOwner struct {
	Usage KeyUsage
	Value string
}

// Key contains an access token and its owner.
type Key struct {
	ClientID  null.String `json:"clientId" db:"client_id"` // Input.
	CreatedBy string      `json:"createdBy" db:"created_by"`
	Token     string      `json:"token" db:"token"`
	Usage     KeyUsage    `json:"usage" db:"usage_type"`
}
