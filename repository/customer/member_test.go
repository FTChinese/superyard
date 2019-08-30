package customer

import (
	"github.com/FTChinese/go-rest/chrono"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
	"time"
)

func TestEnv_CreateMember(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		m reader.Membership
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create Member",
			args: args{m: test.NewProfile().Membership()},
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

func TestEnv_DeleteMember(t *testing.T) {
	env := Env{DB: test.DBX}

	m := test.NewProfile().Membership()

	if err := env.CreateMember(m); err != nil {
		t.Error(err)
	}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Delete a member",
			args:    args{id: m.ID.String},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DeleteMember(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_RetrieveMember(t *testing.T) {

	env := Env{DB: test.DBX}

	m := test.NewProfile().Membership()

	if err := env.CreateMember(m); err != nil {
		t.Error(err)
	}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve member by id",
			args:    args{id: m.ID.String},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveMember(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Membership: %+v", got)
		})
	}
}

func TestEnv_UpdateMember(t *testing.T) {
	env := Env{DB: test.DBX}

	m := test.NewProfile().Membership()

	m.ExpireDate = chrono.DateFrom(time.Now().AddDate(2, 0, 0))

	if err := env.CreateMember(m); err != nil {
		t.Error(err)
	}

	type args struct {
		m reader.Membership
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
