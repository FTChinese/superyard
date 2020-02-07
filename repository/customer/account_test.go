package customer

import (
	"gitlab.com/ftchinese/superyard/test"
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
			got, err := env.RetrieveAccountFtc(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAccountFtc() error = %v, wantErr %v", err, tt.wantErr)
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
			got, err := env.RetrieveAccountWx(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAccountWx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wx Account: %+v", got)
		})
	}
}

func TestEnv_RetrieveMemberFtc(t *testing.T) {
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
			name:    "Retrieve FTC Member",
			args:    args{ftcID: test.MyProfile.FtcID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveMemberFtc(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveMemberFtc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("FTC Membership: %+v", got)
		})
	}
}

func TestEnv_RetrieveMemberWx(t *testing.T) {
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
			name:    "Retrieve Wechat Member",
			args:    args{unionID: test.MyProfile.UnionID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveMemberWx(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveMemberWx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wechat Membership: %+v", got)
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
