package readers

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestEnv_CreateSandboxUser(t *testing.T) {
	p := test.NewPersona()

	t.Logf("ID: %s", p.FtcID)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		account reader.FtcAccount
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create sandbox user",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				account: p.FtcAccount(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.CreateTestUser(tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("CreateTestUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ListSandboxFtcAccount(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []reader.FtcAccount
		wantErr bool
	}{
		{
			name: "List sandbox users",
			fields: fields{
				db: test.DBX,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.ListTestFtcAccount()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTestFtcAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_sandboxJoinedSchema(t *testing.T) {
	p := test.NewPersona()

	_ = NewEnv(test.DBX).CreateTestUser(p.FtcAccount())

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ftcId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Retrieve a sandbox user",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				ftcId: p.FtcID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.testJoinedSchema(tt.args.ftcId)
			if (err != nil) != tt.wantErr {
				t.Errorf("testJoinedSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_LoadSandboxAccount(t *testing.T) {
	p := test.NewPersona()

	_ = NewEnv(test.DBX).CreateTestUser(p.FtcAccount())

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ftcID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Load sandbox account",
			fields: fields{
				db: test.DBX,
			},
			args:    args{ftcID: p.FtcID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.LoadSandboxAccount(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTestAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_SandboxUserExists(t *testing.T) {
	p := test.NewPersona()

	_ = NewEnv(test.DBX).CreateTestUser(p.FtcAccount())

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Not exists",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				id: test.NewPersona().FtcID,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sandbox user exists",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				id: p.FtcID,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.SandboxUserExists(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SandboxUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SandboxUserExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_ChangePassword(t *testing.T) {
	p := test.NewPersona()

	account := p.FtcAccount()

	t.Logf("Initial password: %s", account.Password)

	_ = NewEnv(test.DBX).CreateTestUser(account)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		u reader.TestPasswordUpdater
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Change password",
			fields: fields{
				db: test.DBX,
			},
			args:    args{u: p.PasswordUpdater()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			t.Logf("New password: %s", tt.args.u.Password)
			if err := env.ChangePassword(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
