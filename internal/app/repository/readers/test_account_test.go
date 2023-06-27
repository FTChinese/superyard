package readers

import (
	"testing"

	"github.com/FTChinese/superyard/internal/pkg/sandbox"
	"github.com/FTChinese/superyard/pkg/db"
	"go.uber.org/zap/zaptest"
)

func TestEnv_CreateTestUser(t *testing.T) {
	env := New(db.MockGormSQL(), zaptest.NewLogger(t))

	type args struct {
		account sandbox.TestAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create test user",
			args: args{
				account: sandbox.MockTestAccount(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.CreateTestUser(tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("CreateTestUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestEnv_DeleteTestAccount(t *testing.T) {
// 	env := New(db.MockMySQL(), zaptest.NewLogger(t))

// 	ta := sandbox.MockTestAccount()

// 	_ = env.CreateTestUser(ta)
// 	type args struct {
// 		id string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Delete test account",
// 			args: args{
// 				id: ta.FtcID,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if err := env.DeleteTestAccount(tt.args.id); (err != nil) != tt.wantErr {
// 				t.Errorf("DeleteTestAccount() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func TestEnv_LoadSandboxAccount(t *testing.T) {

	env := New(db.MockGormSQL(), zaptest.NewLogger(t))

	ta := sandbox.MockTestAccount()

	_ = env.CreateTestUser(ta)

	type args struct {
		ftcID string
	}
	tests := []struct {
		name    string
		args    args
		want    sandbox.TestAccount
		wantErr bool
	}{
		{
			name: "Load sandbox account",
			args: args{
				ftcID: ta.FtcID,
			},
			want:    sandbox.TestAccount{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadSandboxAccount(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSandboxAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadSandboxAccount() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_ChangePassword(t *testing.T) {
	env := New(db.MockGormSQL(), zaptest.NewLogger(t))

	ta := sandbox.MockTestAccount()

	_ = env.CreateTestUser(ta)

	type args struct {
		a sandbox.TestAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Change password",
			args: args{
				a: ta,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.ChangePassword(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
