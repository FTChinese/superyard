package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/test"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_ListWxLoginHistory(t *testing.T) {
	env := New(db.MustNewMyDBs(false), zaptest.NewLogger(t))

	type args struct {
		unionID string
		p       gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Wechat Login History",
			args:    args{unionID: test.MyProfile.UnionID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListWxLoginHistory(tt.args.unionID, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListWxLoginHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wx login history: %+v", got)
		})
	}
}
