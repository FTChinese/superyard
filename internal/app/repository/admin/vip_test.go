package admin

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_FtcAccount(t *testing.T) {
	p := test.NewPersona().SetVIP()

	test.NewRepo().MustCreateReader(p.FtcAccount())

	env := NewEnv(db.MustNewMyDBs(false))
	type args struct {
		ftcID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ftc account by uuid",
			args: args{
				ftcID: p.FtcID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.FtcAccount(tt.args.ftcID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_UpdateVIP(t *testing.T) {
	p := test.NewPersona().SetVIP()
	test.NewRepo().MustCreateReader(p.FtcAccount())

	env := NewEnv(db.MustNewMyDBs(false))
	type args struct {
		a reader.FtcAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set vip",
			args: args{
				a: p.FtcAccount(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.UpdateVIP(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("UpdateVIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
