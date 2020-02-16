package readers

import (
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_ListEmailLoginHistory(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		ftcID string
		p     util.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Email Login History",
			args: args{
				ftcID: test.MyProfile.FtcID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListEmailLoginHistory(tt.args.ftcID, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListEmailLoginHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Email Login History: %+v", got)
		})
	}
}

func TestEnv_ListWxLoginHistory(t *testing.T) {
	env := Env{DB: test.DBX}

	type args struct {
		unionID string
		p       util.Pagination
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
