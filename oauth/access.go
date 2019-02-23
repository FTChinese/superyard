package oauth

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
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
	ID          int64       `json:"id"`
	Token       string      `json:"token"`
	Description null.String `json:"description"` // Optional user input data. Max 256
	CreatedAt   chrono.Time `json:"createdAt"`
	UpdatedAt   chrono.Time `json:"updatedAt"`
	LastUsedAt  chrono.Time `json:"lastUsedAt"`
}

// NewAccess creates a new access token instance with token generated.
func NewAccess() (Access, error) {
	acc := Access{}
	t, err := NewToken()
	if err != nil {
		return acc, err
	}

	acc.Token = t

	return acc, nil
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
