package readers

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_MemberByCompoundID(t *testing.T) {
	p := test.NewPersona()

	test.NewRepo().MustCreateMembership(p.Membership())
	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))
	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Retrieve by compound id",
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

func TestEnv_DeleteFtcMember(t *testing.T) {
	p := test.NewPersona()
	test.NewRepo().MustCreateMembership(p.Membership())

	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))

	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete ftc member",
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
