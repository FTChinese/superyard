package readers

import (
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/test"
)

func TestEnv_SearchFtcAccount(t *testing.T) {

	r := test.NewPersona()

	test.NewRepo().MustCreateReader(r)

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
			name:    "Search account by email",
			args:    args{email: r.Email},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.SearchFtcAccount(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchFtcAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Found account: %+v", got)
		})
	}
}

func TestEnv_SearchWxAccounts(t *testing.T) {
	r := test.NewPersona()

	test.NewRepo().MustCreateWxInfo(r.WxInfo())

	env := Env{DB: test.DBX}

	type args struct {
		nickname string
		p        util.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []reader.BaseAccount
		wantErr bool
	}{
		{
			name: "Search Wx account",
			args: args{
				nickname: r.Nickname,
				p: util.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
		},
		{
			name: "No result",
			args: args{
				nickname: gofakeit.Username(),
				p: util.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.SearchWxAccounts(tt.args.nickname, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchWxAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wx readers: %+v", got)
		})
	}
}

func TestEnv_findWxIDs(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		nickname string
		p        util.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Find no result",
			args: args{
				nickname: gofakeit.Username(),
				p:        util.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.findWxIDs(tt.args.nickname, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.findWxIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Uninon IDs %+v", got)
		})
	}
}

func TestEnv_retrieveWxAccounts(t *testing.T) {
	env := Env{DB: test.DBX}

	type args struct {
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "No result",
			args: args{
				ids: []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.retrieveWxAccounts(tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.retrieveWxAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wechat acocunt %+v", got)
		})
	}
}
