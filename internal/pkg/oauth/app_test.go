package oauth

import (
	"testing"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/asaskevich/govalidator"
	"github.com/guregu/null"
)

func TestApp_Validate(t *testing.T) {
	id, _ := conv.RandomHexBin(10)
	secret, _ := conv.RandomHexBin(32)
	app := App{
		BaseApp: BaseApp{
			Name:        "superyard",
			Slug:        "superyard",
			RepoURL:     "https://github.com/FTChinese/superyard-go",
			Description: null.String{},
			HomeURL:     null.String{},
			CallbackURL: null.String{},
		},
		ClientID:     id,
		ClientSecret: secret,
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
