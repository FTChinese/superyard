package oauth

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/util"
	"strings"
)

// APIKey is an OAuth 2.0 access token used by an app or person to access ftc api
type Access struct {
	ID          int    `json:"id"`
	token       string `json:"token"`       // required but auto generated, 20 bytes
	Description string `json:"description"` // Required, max 255 chars
	MyftID      string `json:"myftId"`      // optional, ftc account associated with this access token.
	CreateAt    string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	LastUsedAt  string `json:"lastUsedAt"`
	CreatedBy   string `json:"createdBy"`  // optional, for personal access token.
	OwnedByApp  string `json:"ownedByApp"` // optional, for client access token.
}

// Sanitize removes leading and trailing spaces
func (a *Access) Sanitize() {
	a.Description = strings.TrimSpace(a.Description)
	a.MyftID = strings.TrimSpace(a.MyftID)
	a.OwnedByApp = strings.TrimSpace(a.OwnedByApp)
}

// Validate checks max length of each fields
func (a Access) Validate() *view.Reason {
	return util.OptionalMaxLen(a.Description, 255, "description")
}

func (a Access) GetToken() string {
	return a.token
}

func NewAccess() (Access, error)  {
	token, err := gorest.RandomHex(20)

	if err != nil {
		return Access{}, err
	}

	return Access{token: token}, nil
}