package oauth

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/validator"
	"strings"
)

type AppRemover struct {
	ClientID string `db:"client_id"`
	OwnedBy  string `json:"ownedBy" db:"owned_by"`
}

// BaseApp represents the input body of a request when creating an app.
type BaseApp struct {
	Name        string      `json:"name" db:"app_name" valid:"required,length(1|256)"`        // required, max 256 chars. Can be updated.
	Slug        string      `json:"slug" db:"slug_name" valid:"required,length(1|256)"`       // required, unique, max 255 chars
	RepoURL     string      `json:"repoUrl" db:"repo_url" valid:"required,url,length(1|256)"` // required, 256 chars. Can be updated.
	Description null.String `json:"description" db:"description" valid:"-"`                   // optional, 512 chars. Can be updated.
	HomeURL     null.String `json:"homeUrl" db:"home_url" valid:"-"`                          // optional, 256 chars. Can be updated.
	CallbackURL null.String `json:"callbackUrl" db:"callback_url" valid:"-"`
}

// Sanitize removes leading and trailing spaces
func (a *BaseApp) Sanitize() {
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
func (a BaseApp) Validate() *render.ValidationError {
	ve := validator.New("name").Required().Max(256).Validate(a.Name)
	if ve != nil {
		return ve
	}

	ve = validator.New("slug").Required().Max(256).Validate(a.Slug)
	if ve != nil {
		return ve
	}

	ve = validator.New("repoUrl").Required().Max(256).Validate(a.RepoURL)
	if ve != nil {
		return ve
	}

	ve = validator.New("description").Max(512).Validate(a.Description.String)
	if ve != nil {
		return ve
	}

	return validator.New("homeUrl").Max(256).Validate(a.HomeURL.String)
}

// App represents an application that needs to access ftc api
type App struct {
	BaseApp
	ClientID     string      `json:"clientId" db:"client_id" valid:"-"`         // required, 10 bytes. Immutable once created.
	ClientSecret string      `json:"clientSecret" db:"client_secret" valid:"-"` // required, 32 bytes. Immutable once created.
	IsActive     bool        `json:"isActive" db:"is_active" valid:"-"`
	CreatedAt    chrono.Time `json:"createdAt" db:"created_at" valid:"-"`
	UpdatedAt    chrono.Time `json:"updatedAt" db:"updated_at" valid:"-"`
	OwnedBy      string      `json:"ownedBy" db:"owned_by" valid:"-"`
}

func NewApp(base BaseApp, owner string) (App, error) {
	clientID, err := gorest.RandomHex(10)
	if err != nil {
		return App{}, err
	}

	clientSecret, err := gorest.RandomHex(32)
	if err != nil {
		return App{}, err
	}

	return App{
		BaseApp:      base,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		IsActive:     false,
		CreatedAt:    chrono.Time{},
		UpdatedAt:    chrono.Time{},
		OwnedBy:      owner,
	}, nil
}
