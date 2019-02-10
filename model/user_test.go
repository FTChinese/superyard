package model

import (
	"database/sql"
	"github.com/FTChinese/go-rest/enum"
	"log"
	"testing"

	"github.com/guregu/null"
)

func TestUserEnv_LoadAccount(t *testing.T) {
	m := newMockUser().withUnionID()
	u := m.createUser()
	m.createWxUser()

	order := m.order(stdPlan, enum.LoginMethodEmail)
	m.createMember(order)

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
			name:    "Load Account",
			fields:  fields{DB: db},
			args:    args{userID: u.UserID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := UserEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LoadAccount(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserEnv.LoadAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestUserEnv_ListOrders(t *testing.T) {
	mUser := newMockUser().withUnionID()
	mUser.createUser()

	o1 := mUser.order(stdPlan, enum.LoginMethodEmail)
	o2 := mUser.order(stdPlan, enum.LoginMethodWx)

	mUser.createOrder(o1)
	mUser.createOrder(o2)

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
