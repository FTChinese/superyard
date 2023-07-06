package registry

import (
	"testing"

	"github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/brianvoe/gofakeit/v5"
)

func mockNewToken() oauth.Access {
	a, err := oauth.NewAccess(oauth.BaseAccess{
		ClientID: conv.MustRandomHexBin(10),
	}, gofakeit.Username())

	if err != nil {
		panic(err)
	}

	return a
}

func mockNewAppToken(app oauth.App) oauth.Access {
	a, err := oauth.NewAccess(oauth.BaseAccess{
		ClientID: app.ClientID,
	}, gofakeit.Username())

	if err != nil {
		panic(err)
	}

	return a
}

func mockNewPersionalToken(creator string) oauth.Access {
	a, err := oauth.NewAccess(oauth.BaseAccess{}, creator)

	if err != nil {
		panic(err)
	}

	return a
}

func TestEnv_CreateToken(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a := mockNewToken()

	type args struct {
		acc oauth.Access
	}
	tests := []struct {
		name    string
		args    args
		want    oauth.Access
		wantErr bool
	}{
		{
			name: "create token",
			args: args{
				acc: a,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.CreateToken(tt.args.acc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_RetrieveToken(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a := mockNewToken()

	a, _ = env.CreateToken(a)

	type args struct {
		id    int64
		owner string
	}
	tests := []struct {
		name    string
		args    args
		want    oauth.Access
		wantErr bool
	}{
		{
			name: "retrieve token",
			args: args{
				id:    a.ID,
				owner: a.CreatedBy,
			},
			want:    a,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveToken(tt.args.id, tt.args.owner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.RetrieveToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("Env.RetrieveToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_ListAppTokens(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	app := mockNewApp()

	env.CreateToken(mockNewAppToken(app))
	env.CreateToken(mockNewAppToken(app))
	env.CreateToken(mockNewAppToken(app))

	type args struct {
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    []oauth.Access
		wantErr bool
	}{
		{
			name: "list app token",
			args: args{
				clientID: app.ClientID.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListAppTokens(tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.ListAppTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%d", len(got))
		})
	}
}

func TestEnv_ListPersonalKeys(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	creator := gofakeit.Username()

	env.CreateToken(mockNewPersionalToken(creator))
	env.CreateToken(mockNewPersionalToken(creator))
	env.CreateToken(mockNewPersionalToken(creator))

	type args struct {
		owner string
	}
	tests := []struct {
		name    string
		args    args
		want    []oauth.Access
		wantErr bool
	}{
		{
			name: "list personal key",
			args: args{
				owner: creator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListPersonalKeys(tt.args.owner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.ListPersonalKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%d", len(got))
		})
	}
}

func TestEnv_RemoveKey(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a := mockNewToken()

	a, _ = env.CreateToken(a)

	a = a.Remove()

	type args struct {
		k oauth.Access
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove key",
			args: args{
				k: a,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.RemoveKey(tt.args.k); (err != nil) != tt.wantErr {
				t.Errorf("Env.RemoveKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
