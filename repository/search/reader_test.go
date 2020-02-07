package search

import (
	"github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_SearchFtcUser(t *testing.T) {
	env := Env{DB: test.DBX}

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Search FTC User",
			args:    args{email: test.MyProfile.Email},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.SearchFtcUser(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchFtcUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Ftc user: %+v", got)
		})
	}
}

func TestEnv_SearchWxUser(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		nickname string
		p        builder.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Search wechat user",
			args: args{
				nickname: test.MyProfile.Nickname,
				p:        builder.NewPagination(1, 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.SearchWxUser(tt.args.nickname, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchWxUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wechat user: %+v", got)
		})
	}
}
