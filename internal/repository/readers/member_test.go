package readers

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_MemberByCompoundID(t *testing.T) {
	p := test.NewPersona()

	test.NewRepo().MustCreateMembership(p.Membership())
	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))
	type fields struct {
		db     *sqlx.DB
		logger *zap.Logger
	}
	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Retrieve by compound id",
			fields: fields{
				db:     test.DBX,
				logger: zaptest.NewLogger(t),
			},
			args: args{
				compoundID: p.FtcID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.MemberByCompoundID(tt.args.compoundID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_CreateFtcMember(t *testing.T) {
	p := test.NewPersona()

	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))

	type fields struct {
		db     *sqlx.DB
		logger *zap.Logger
	}
	type args struct {
		input subs.FtcSubsCreationInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create ftc member",
			fields: fields{
				db:     test.DBX,
				logger: zaptest.NewLogger(t),
			},
			args: args{
				input: p.FtcSubsCreationInput(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.CreateFtcMember(tt.args.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_UpdateFtcMember(t *testing.T) {
	p := test.NewPersona()
	test.NewRepo().MustCreateMembership(p.Membership())

	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))

	type fields struct {
		db     *sqlx.DB
		logger *zap.Logger
	}
	type args struct {
		compoundID string
		input      subs.FtcSubsUpdateInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update ftc member",
			fields: fields{
				db:     test.DBX,
				logger: zaptest.NewLogger(t),
			},
			args: args{
				compoundID: p.FtcID,
				input:      p.FtcSubsUpdateInput(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.UpdateFtcMember(tt.args.compoundID, tt.args.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_DeleteFtcMember(t *testing.T) {
	p := test.NewPersona()
	test.NewRepo().MustCreateMembership(p.Membership())

	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))

	type fields struct {
		db     *sqlx.DB
		logger *zap.Logger
	}
	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete ftc member",
			fields: fields{
				db:     test.DBX,
				logger: zaptest.NewLogger(t),
			},
			args: args{
				compoundID: p.FtcID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.DeleteFtcMember(tt.args.compoundID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}
