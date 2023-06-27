package oauth

import (
	"strings"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
)

// AppRemoved is used to identify an app to be removed.
type AppRemover struct {
	ClientID string `db:"client_id"`
	OwnedBy  string `json:"ownedBy" db:"owned_by"`
}

// BaseApp represents the input body of a request when creating an app.
type BaseApp struct {
	Name        string      `json:"name" db:"app_name"`           // required, max 256 chars. Can be updated.
	Slug        string      `json:"slug" db:"slug_name"`          // required, unique, max 255 chars
	RepoURL     string      `json:"repoUrl" db:"repo_url"`        // required, 256 chars. Can be updated.
	Description null.String `json:"description" db:"description"` // optional, 512 chars. Can be updated.
	HomeURL     null.String `json:"homeUrl" db:"home_url"`        // optional, 256 chars. Can be updated.
	CallbackURL null.String `json:"callbackUrl" db:"callback_url"`
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
	ve := validator.New("name").
		Required().
		MaxLen(64).
		Validate(a.Name)
	if ve != nil {
		return ve
	}

	ve = validator.New("slug").
		Required().
		MaxLen(64).
		Validate(a.Slug)
	if ve != nil {
		return ve
	}

	ve = validator.New("repoUrl").
		Required().
		MaxLen(256).
		Validate(a.RepoURL)
	if ve != nil {
		return ve
	}

	ve = validator.New("description").
		MaxLen(512).
		Validate(a.Description.String)
	if ve != nil {
		return ve
	}

	return validator.New("homeUrl").
		MaxLen(256).
		Validate(a.HomeURL.String)
}

// App represents an application that needs to access ftc api
type App struct {
	BaseApp
	ClientID     conv.HexBin `json:"clientId" db:"client_id" gorm:"column:client_id"`             // required, 10 bytes. Immutable once created.
	ClientSecret conv.HexBin `json:"clientSecret" db:"client_secret" gorm:"column:client_secret"` // required, 32 bytes. Immutable once created.
	IsActive     bool        `json:"isActive" db:"is_active" gorm:"column:is_active"`
	CreatedAt    chrono.Time `json:"createdAt" db:"created_at" gorm:"column:created_at"`
	UpdatedAt    chrono.Time `json:"updatedAt" db:"updated_at" gorm:"column:updated_at"`
	OwnedBy      string      `json:"ownedBy" db:"owned_by" gorm:"column:owned_by"`
}

func (App) TableName() string {
	return "oauth.app_registry"
}

func NewApp(base BaseApp, owner string) (App, error) {
	clientID, err := conv.RandomHexBin(10)
	if err != nil {
		return App{}, err
	}

	clientSecret, err := conv.RandomHexBin(32)
	if err != nil {
		return App{}, err
	}

	return App{
		BaseApp:      base,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		IsActive:     true,
		CreatedAt:    chrono.Time{},
		UpdatedAt:    chrono.Time{},
		OwnedBy:      owner,
	}, nil
}

type AppList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []App `json:"data"`
	Err  error `json:"-"`
}
