package oauth

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/util"
	"strings"
)

// App represents an application that needs to access ftc api
type App struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`         // required, max 255 chars. Can be updated.
	Slug         string      `json:"slug"`         // required, unique, max 255 chars
	ClientID     string      `json:"clientId"`     // required, 10 bytes. Immutable once created.
	ClientSecret string      `json:"clientSecret"` // required, 32 bytes. Immutable once created.
	RepoURL      string      `json:"repoUrl"`      // required, 255 chars. Can be updated.
	Description  null.String `json:"description"`  // optional, 511 chars. Can be updated.
	HomeURL      null.String `json:"homeUrl"`      // optional, 255 chars. Can be updated.
	IsActive     bool        `json:"isActive"`
	CreatedAt    chrono.Time `json:"createdAt"`
	UpdatedAt    chrono.Time `json:"updatedAt"`
	OwnedBy      string      `json:"ownedBy"`
}

// Sanitize removes leading and trailing spaces
func (a *App) Sanitize() {
	a.Name = strings.TrimSpace(a.Name)
	a.Slug = strings.TrimSpace(a.Slug)
	a.RepoURL = strings.TrimSpace(a.RepoURL)

	if a.Description.Valid {
		 a.Description.String = strings.TrimSpace(a.Description.String)
	}

	if a.HomeURL.Valid {
		a.HomeURL.String = strings.TrimSpace(a.HomeURL.String)
	}
}

// Validate performs validation on incoming app.
func (a App) Validate() *view.Reason {
	if r := util.RequireNotEmptyWithMax(a.Name, 255, "name"); r != nil {
		return r
	}

	if r := util.RequireNotEmptyWithMax(a.Slug, 255, "slug"); r != nil {
		return r
	}

	if r := util.RequireNotEmptyWithMax(a.RepoURL, 255, "repoUrl"); r != nil {
		return r
	}

	if r := util.OptionalMaxLen(a.Description.String, 500, "description"); r != nil {
		return r
	}

	return util.OptionalMaxLen(a.HomeURL.String, 120, "homeUrl")
}

// GenCredentials generates SlugName and ClientSecret.
func (a *App) GenCredentials() error {
	clientID, err := gorest.RandomHex(10)
	if err != nil {
		return err
	}
	
	a.ClientID = clientID

	clientSecret, err := gorest.RandomHex(32)
	if err != nil {
		return err
	}
	
	a.ClientSecret = clientSecret

	return nil
}

// Ownership is used to transfer an app's ownership
type Ownership struct {
	SlugName string
	NewOwner string
	OldOwner string
}
