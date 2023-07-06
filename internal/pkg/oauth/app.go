package oauth

import (
	"strings"

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
	Name        string      `json:"name" gorm:"column:app_name"`            // required, max 256 chars. Can be updated.
	Slug        string      `json:"slug" gorm:"column:slug_name"`           // required, unique, max 255 chars
	RepoURL     string      `json:"repoUrl" gorm:"column:repo_url"`         // required, 256 chars. Can be updated.
	Description null.String `json:"description"  gorm:"column:description"` // optional, 512 chars. Can be updated.
	HomeURL     null.String `json:"homeUrl" gorm:"column:homepage_url"`     // optional, 256 chars. Can be updated.
	CallbackURL null.String `json:"callbackUrl"  gorm:"column:callback_url"`
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
	ID int64 `json:"id" gorm:"primaryKey"`
	BaseApp
	ClientID     conv.HexBin `json:"clientId" gorm:"column:client_id"`         // required, 10 bytes. Immutable once created.
	ClientSecret conv.HexBin `json:"clientSecret" gorm:"column:client_secret"` // required, 32 bytes. Immutable once created.
	IsActive     bool        `json:"isActive" gorm:"column:is_active"`
	CreatedAt    chrono.Time `json:"createdAt" gorm:"column:created_utc"`
	UpdatedAt    chrono.Time `json:"updatedAt" gorm:"column:updated_utc"`
	OwnedBy      string      `json:"ownedBy" gorm:"column:owned_by"`
}

func (App) TableName() string {
	return "oauth.app_registry"
}

func (a App) Update(input BaseApp) App {
	a.BaseApp = input
	a.UpdatedAt = chrono.TimeNow()
	return a
}

func (a App) Remove() App {
	a.IsActive = false
	a.UpdatedAt = chrono.TimeNow()
	return a
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
		CreatedAt:    chrono.TimeNow(),
		UpdatedAt:    chrono.TimeNow(),
		OwnedBy:      owner,
	}, nil
}
