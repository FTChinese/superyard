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
		account reader.SandboxUser
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
				account: p.Reader(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.CreateSandboxUser(tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("CreateSandboxUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ListSandboxUsers(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []reader.SandboxUser
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
			got, err := env.ListSandboxUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSandboxUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}

func TestEnv_sandboxUserInfo(t *testing.T) {
	p := test.NewPersona()

	_ = NewEnv(test.DBX).CreateSandboxUser(p.Reader())

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
			got, err := env.sandboxUserInfo(tt.args.ftcId)
			if (err != nil) != tt.wantErr {
				t.Errorf("sandboxUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}

func TestEnv_LoadSandboxAccount(t *testing.T) {
	p := test.NewPersona()

	_ = NewEnv(test.DBX).CreateSandboxUser(p.Reader())

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
				t.Errorf("LoadSandboxAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}

func TestEnv_SandboxUserExists(t *testing.T) {
	p := test.NewPersona()

	_ = NewEnv(test.DBX).CreateSandboxUser(p.Reader())

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

	u := p.Reader()

	t.Logf("Initial password: %s", u.Password)

	u.Password = faker.SimplePassword()

	t.Logf("Changed password to %s", u.Password)

	_ = NewEnv(test.DBX).CreateSandboxUser(p.Reader())

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		u reader.SandboxUser
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
			args:    args{u: u},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.ChangePassword(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
