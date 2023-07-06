package readers

import (
	"testing"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"go.uber.org/zap/zaptest"
)

func mockNewWxProfile() reader.WxProfile {
	id, _ := gorest.RandomBase64(10)

	return reader.WxProfile{
		UnionID:   id,
		Nickname:  null.NewString(gofakeit.Username(), true),
		AvatarURL: null.NewString(gofakeit.URL(), true),
		Gender:    enum.GenderFemale,
		Country:   null.NewString(gofakeit.Country(), true),
		Province:  null.NewString(gofakeit.State(), true),
		City:      null.NewString(gofakeit.City(), true),
		CreatedAt: chrono.TimeNow(),
		UpdatedAt: chrono.TimeNow(),
	}
}

func (env Env) createWxProfile(p reader.WxProfile) {
	err := env.gormDBs.Write.Create(&p).Error
	if err != nil {
		panic(err)
	}
}

func TestEnv_RetrieveWxProfile(t *testing.T) {
	env := New(db.MockGormSQL(), zaptest.NewLogger(t))

	p := mockNewWxProfile()

	env.createWxProfile(p)

	type args struct {
		unionID string
	}
	tests := []struct {
		name    string
		args    args
		want    reader.WxProfile
		wantErr bool
	}{
		{
			name: "retrieve wechat profile",
			args: args{
				unionID: p.UnionID,
			},
			want:    p,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveWxProfile(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.RetrieveWxProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.UnionID != tt.want.UnionID {
				t.Errorf("Env.RetrieveWxProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
