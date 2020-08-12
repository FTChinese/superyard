package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"testing"
)

func newBanner() paywall.Banner {

	input := paywall.BannerInput{
		Heading:    "",
		CoverURL:   null.String{},
		SubHeading: null.String{},
		Content:    null.String{},
	}

	return paywall.NewBanner(input, "weiguo.ni")
}

func TestEnv_CreateBanner(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		b paywall.Banner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.CreateBanner(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("CreateBanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
