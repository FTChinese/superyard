package auth

import (
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
