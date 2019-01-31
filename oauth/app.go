package oauth

import (
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/util"
	"strings"
)

// App represents an application that needs to access ftc api
type App struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`         // required, max 255 chars.
	Slug         string `json:"slug"`         // required, unique, max 255 chars
	ClientID     string `json:"clientId"`     // required, 10 bytes
	ClientSecret string `json:"clientSecret"` // required, 32 bytes
	RepoURL      string `json:"repoUrl"`      // required, 255 chars
	Description  string `json:"description"`  // optional, 511 chars
	HomeURL      string `json:"homeUrl"`      // optional, 255 chars
	IsActive     bool   `json:"isActive"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	OwnedBy      string `json:"ownedBy"`
}

// Sanitize removes leading and trailing spaces
func (a *App) Sanitize() {
	a.Name = strings.TrimSpace(a.Name)
	a.Slug = strings.TrimSpace(a.Slug)
	a.RepoURL = strings.TrimSpace(a.RepoURL)
	a.Description = strings.TrimSpace(a.Description)
	a.HomeURL = strings.TrimSpace(a.HomeURL)
}

// Validate performas validation on incoming app.
func (a *App) Validate() *view.Reason {
	if r := util.RequireNotEmptyWithMax(a.Name, 255, "name"); r != nil {
		return r
	}

	if r := util.RequireNotEmptyWithMax(a.Slug, 255, "slug"); r != nil {
		return r
	}

	if r := util.RequireNotEmptyWithMax(a.RepoURL, 255, "repoUrl"); r != nil {
		return r
	}

	if r := util.OptionalMaxLen(a.Description, 500, "description"); r != nil {
		return r
	}

	return util.OptionalMaxLen(a.HomeURL, 120, "homeUrl")
}

// Ownership is used to transfer an app's ownership
type Ownership struct {
	SlugName string
	NewOwner string
	OldOwner string
}
