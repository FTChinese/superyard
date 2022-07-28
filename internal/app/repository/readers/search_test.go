package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/db"
	"go.uber.org/zap/zaptest"
	"testing"

	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/test"
	"github.com/brianvoe/gofakeit/v5"
)

func TestEnv_SearchWxAccounts(t *testing.T) {
	r := test.NewPersona()

	test.NewRepo().MustCreateWxInfo(r.WxInfo())

	env := New(db.MustNewMyDBs(false), zaptest.NewLogger(t))

	type args struct {
		nickname string
		p        gorest.Pagination
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
				p: gorest.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
		},
		{
			name: "No result",
			args: args{
				nickname: gofakeit.Username(),
				p: gorest.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.SearchJoinedAccountWxName(tt.args.nickname, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchJoinedAccountWxName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wx readers: %+v", got)
		})
	}
}
