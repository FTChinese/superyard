package model

import (
	"database/sql"
	"github.com/FTChinese/go-rest/enum"
	"log"
	"reflect"
	"testing"

	"gitlab.com/ftchinese/backyard-api/user"
)

func TestSearchEnv_FindUserByEmail(t *testing.T) {
	m := newMockUser()
	m.createUser()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name:    "Find User By Email",
			fields:  fields{DB: db},
			args:    args{email: m.email},
			want:    m.user(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := SearchEnv{
				DB: tt.fields.DB,
			}
			got, err := env.FindUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchEnv.FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchEnv.FindUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchEnv_FindOrder(t *testing.T) {
	mUser := newMockUser()
	order := mUser.order(stdPlan, enum.LoginMethodEmail)

	mUser.createOrder(order)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		orderID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Find Order",
			fields:  fields{DB: db},
			args:    args{orderID: order.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := SearchEnv{
				DB: tt.fields.DB,
			}
			got, err := env.FindOrder(tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchEnv.FindOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("%+v", got)
		})
	}
}
