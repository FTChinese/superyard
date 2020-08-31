package readers

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/test"
	"github.com/guregu/null"
	"testing"
	"time"
)

func TestMemberTx_RetrieveMember(t *testing.T) {

	p := test.NewPersona()

	test.NewRepo().MustCreateMembership(p.Membership())

	tx, _ := NewEnv(test.DBX).BeginMemberTx()

	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Retrieve membership",
			args: args{
				compoundID: p.FtcID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tx.RetrieveMember(tt.args.compoundID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}

	if err := tx.Commit(); err != nil {
		t.Error(err)
	}
}

func TestMemberTx_CreateMember(t *testing.T) {

	p := test.NewPersona()

	tx, _ := NewEnv(test.DBX).BeginMemberTx()

	type args struct {
		m reader.Membership
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create member",
			args: args{
				m: p.Membership(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tx.CreateMember(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("CreateMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	if err := tx.Commit(); err != nil {
		t.Error(err)
	}
}

func TestMemberTx_UpdateMember(t *testing.T) {
	p := test.NewPersona()
	m := p.Membership()

	test.NewRepo().MustCreateMembership(m)

	tx, _ := NewEnv(test.DBX).BeginMemberTx()

	type args struct {
		m reader.Membership
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update member",
			args: args{
				m: m.Update(reader.MemberInput{
					ExpireDate: chrono.DateFrom(time.Now().AddDate(2, 0, 0)),
					PayMethod:  enum.PayMethodWx,
					FtcPlanID:  null.StringFrom(test.PlanPrm.ID),
				}, test.PlanPrm.Plan),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tx.UpdateMember(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	if err := tx.Commit(); err != nil {
		t.Error(err)
	}
}

func TestMemberTx_DeleteMember(t *testing.T) {
	p := test.NewPersona()
	m := p.Membership()

	test.NewRepo().MustCreateMembership(m)

	tx, _ := NewEnv(test.DBX).BeginMemberTx()

	type args struct {
		compoundID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete a member",
			args: args{
				compoundID: m.CompoundID.String,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tx.DeleteMember(tt.args.compoundID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	if err := tx.Commit(); err != nil {
		t.Error(err)
	}
}
