package readers

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
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

	t.Logf("Compound id: %s", m.CompoundID)

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
			name:    "Update Member",
			args:    args{m: m},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateMember(tt.args.m, "weiguo.ni"); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_FindMemberForOrder(t *testing.T) {
	p := test.NewPersona().
		SetAccountKind(reader.AccountKindLinked)
	repo := test.NewRepo()

	repo.MustCreateMembership(p.Membership())

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ftcOrUnionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "By ftc id",
			fields: fields{DB: test.DBX},
			args: args{
				ftcOrUnionID: p.FtcID,
			},
			wantErr: false,
		},
		{
			name: "By union id",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				ftcOrUnionID: p.UnionID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.FindMemberForOrder(tt.args.ftcOrUnionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMemberForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotZero(t, got)
		})
	}
}
