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
	ID         int         `json:"id"`
	Token      string      `json:"token"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdatedAt  chrono.Time `json:"updatedAt"`
	LastUsedAt chrono.Time `json:"lastUsedAt"`
}

// PersonalAccess is an api key used by a person.
type PersonalAccess struct {
	Access
	Description null.String `json:"description"` // Optional user input data. Max 256
	MyftEmail   null.String `json:"myftEmail"`   // optional user input data. The ftc account associated with this access Token.
	CreatedBy   null.String `json:"createdBy"`   // optional, for personal access Token.
}
