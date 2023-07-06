package auth

import (
	"reflect"
	"testing"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/brianvoe/gofakeit/v5"
)

type mockAccount struct {
	user.Account
	Password    string      `gorm:"column:password"`
	CreatedAt   chrono.Time `gorm:"column:creatdate"`
	LastLoginAt chrono.Time `gorm:"column:lastlogin"`
}

func (a mockAccount) hashPassword() mockAccount {
	a.Password = conv.NewMD5Sum(a.Password).String()

	return a
}

func mockNewAccount() mockAccount {
	faker.SeedGoFake()

	return mockAccount{
		Account: user.Account{
			UserName:    gofakeit.Username(),
			Email:       gofakeit.Email(),
			DisplayName: gofakeit.FirstName(),
		},
		Password:    faker.SimplePassword(),
		CreatedAt:   chrono.TimeNow(),
		LastLoginAt: chrono.TimeNow(),
	}
}

// Used for mocking only.
func (env Env) createdAccount(a mockAccount) (mockAccount, error) {
	hashed := a.hashPassword()
	err := env.gormDBs.Write.Create(&hashed).Error
	if err != nil {
		return mockAccount{}, err
	}

	a.ID = hashed.ID
	return a, nil
}

func TestEnv_createdAccount(t *testing.T) {

	env := NewEnv(db.MockGormSQL())
	a := mockNewAccount()

	type args struct {
		a mockAccount
	}
	tests := []struct {
		name    string
		args    args
		want    mockAccount
		wantErr bool
	}{
		{
			name: "Create account",
			args: args{
				a: a,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.createdAccount(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.createdAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_AccountByID(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a, _ := env.createdAccount(mockNewAccount())

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    user.Account
		wantErr bool
	}{
		{
			name: "account by id",
			args: args{
				id: a.ID,
			},
			want:    a.Account,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.AccountByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.AccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.AccountByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_SetEmail(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a, _ := env.createdAccount(mockNewAccount())

	type args struct {
		a user.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set email",
			args: args{
				a: a.Account,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SetEmail(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Env.SetEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_UpdateDisplayName(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a, _ := env.createdAccount(mockNewAccount())

	type args struct {
		a user.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "updated display name",
			args: args{
				a: a.Account,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateDisplayName(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Env.UpdateDisplayName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_RetrieveProfile(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a, _ := env.createdAccount(mockNewAccount())

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    user.Profile
		wantErr bool
	}{
		{
			name: "retrieve profile",
			args: args{
				id: a.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveProfile(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.RetrieveProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", got)
		})
	}
}
