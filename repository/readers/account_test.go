package readers

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
			got, err := env.retrieveFTCAccount(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveFTCAccount() error = %v, wantErr %v", err, tt.wantErr)
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
			got, err := env.retrieveWxAccount(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveWxAccount() error = %v, wantErr %v", err, tt.wantErr)
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

			got, err := env.retrieveFtcMember(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveFtcMember() error = %v, wantErr %v", err, tt.wantErr)
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

			got, err := env.retrieveWxMember(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveWxMember() error = %v, wantErr %v", err, tt.wantErr)
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
