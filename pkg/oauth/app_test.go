package oauth

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/asaskevich/govalidator"
	"github.com/guregu/null"
	"testing"
)

func TestApp_Validate(t *testing.T) {
	app := App{
		BaseApp: BaseApp{
			Name:        "superyard",
			Slug:        "superyard",
			RepoURL:     "https://github.com/FTChinese/superyard-go",
			Description: null.String{},
			HomeURL:     null.String{},
			CallbackURL: null.String{},
		},
		ClientID:     "1234567890",
		ClientSecret: "1234567890",
		IsActive:     false,
		CreatedAt:    chrono.Time{},
		UpdatedAt:    chrono.Time{},
		OwnedBy:      "",
	}

	ok, err := govalidator.ValidateStruct(app)

	if err != nil {
		t.Error(err)
	}

	t.Log(ok)
}
