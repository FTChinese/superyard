package readers

import (
	"github.com/FTChinese/superyard/test"
	"testing"
)

func TestEnv_RetrieveAccountFtc(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		ftcID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve FTC Account",
			args:    args{ftcID: test.MyProfile.FtcID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.accountByFtcID(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountByFtcID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("FTC Account: %+v", got)
		})
	}
}

func TestEnv_RetrieveAccountWx(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		unionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve Wechat Account",
			args:    args{unionID: test.MyProfile.UnionID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.accountByWxID(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountByWxID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wx Account: %+v", got)
		})
	}
}

func TestEnv_RetrieveFtcProfile(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		ftcID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve FTC Profile",
			args:    args{ftcID: test.MyProfile.FtcID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveFtcProfile(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveFtcProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("FTC Profile: %+v", got)
		})
	}
}

func TestEnv_RetrieveWxProfile(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		unionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve Wechat Profile",
			args:    args{unionID: test.MyProfile.UnionID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveWxProfile(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveWxProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wx Profile: %+v", got)
		})
	}
}
