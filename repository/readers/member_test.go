package readers

import (
	"github.com/FTChinese/go-rest/chrono"
	"gitlab.com/ftchinese/superyard/pkg/subs"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
	"time"
)

func TestEnv_CreateMember(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		m subs.Membership
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create Member",
			args: args{m: test.NewPersona().Membership()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreateMember(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("CreateMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_UpdateMember(t *testing.T) {
	env := Env{DB: test.DBX}

	m := test.NewPersona().Membership()

	m.ExpireDate = chrono.DateFrom(time.Now().AddDate(2, 0, 0))

	if err := env.CreateMember(m); err != nil {
		t.Error(err)
	}

	type args struct {
		m subs.Membership
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "UpdateProfile Member",
			args:    args{m: m},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateMember(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
