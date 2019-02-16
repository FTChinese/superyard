package model

import (
	"database/sql"
	"gitlab.com/ftchinese/backyard-api/util"
	"log"
	"testing"

	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

func TestUserEnv_LoadAccountByID(t *testing.T) {
	m := newMockUser().withUnionID()
	u := m.createUser()
	m.createWxUser()

	order := m.order(stdPlan, enum.LoginMethodEmail)
	m.createMember(order)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Load Account",
			fields:  fields{DB: db},
			args:    args{id: u.UserID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := UserEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LoadAccountByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserEnv.LoadAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestUserEnv_ListLoginHistory(t *testing.T) {
	m := newMockUser()
	for i := 0; i < 5; i++ {
		m.createLoginHistory()
	}

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "List Login History",
			fields:  fields{DB: db},
			args:    args{userID: m.userID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := UserEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListLoginHistory(tt.args.userID, util.NewPagination(1, 20))
			if (err != nil) != tt.wantErr {
				t.Errorf("UserEnv.ListLoginHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("%+v", got)
		})
	}
}

func TestUserEnv_ListOrders(t *testing.T) {
	mUser := newMockUser().withUnionID()
	mUser.createUser()

	mUser.createOrder(enum.LoginMethodEmail)
	mUser.createOrder(enum.LoginMethodWx)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userID  null.String
		unionID null.String
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Lis Orders",
			fields: fields{DB: db},
			args: args{
				userID:  null.StringFrom(mUser.userID),
				unionID: mUser.unionID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := UserEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListOrders(tt.args.userID, tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserEnv.ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestUserEnv_LoadWxInfo(t *testing.T) {
	m := newMockUser().withUnionID()

	info := m.createWxUser()
	t.Logf("Created wx user: %+v", info)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		unionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Load Wechat User Info",
			fields:  fields{DB: db},
			args:    args{unionID: m.unionID.String},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := UserEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LoadWxInfo(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserEnv.LoadWxInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
		})
	}
}

func TestUserEnv_ListOAuthHistory(t *testing.T) {
	m := newMockUser().withUnionID()
	for i := 0; i < 5; i++ {
		m.createWxAccess()
	}

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		unionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Load Wechat OAuth History",
			fields:  fields{DB: db},
			args:    args{m.unionID.String},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := UserEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListOAuthHistory(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserEnv.ListOAuthHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}

func TestCreateFixedUser(t *testing.T) {
	mocker := newMockUser().withEmail("neefrankie@163.com").withPassword("12345678")

	u := mocker.createUser()
	t.Logf("Created an FTC user: %+v", u)

	for i := 0; i < 5; i++ {
		mocker.createLoginHistory()
	}

	for i := 0; i < 5; i++ {
		mocker.createOrder(enum.LoginMethodEmail)
	}
}
