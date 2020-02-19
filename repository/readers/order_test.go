package readers

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_ConfirmOrder(t *testing.T) {
	env := Env{DB: test.DBX}

	order := test.MyProfile.Order(false)

	if err := env.CreateOrder(order); err != nil {
		t.Error(err)
	}

	orderUp := test.MyProfile.Order(false)
	orderUp.Kind = reader.SubsKindUpgrade
	orderUp.Tier = enum.TierPremium
	if err := env.CreateOrder(orderUp); err != nil {
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
			name:    "Confirm Order",
			args:    args{id: order.ID},
			wantErr: false,
		},
		{
			name:    "Confirm Upgrade Order",
			args:    args{id: orderUp.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.ConfirmOrder(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ConfirmOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_CreateOrder(t *testing.T) {
	env := Env{DB: test.DBX}

	type args struct {
		order reader.Order
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create Order",
			args:    args{order: test.MyProfile.Order(true)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreateOrder(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ListOrders(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		ids reader.AccountID
		p   gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ListStaff Orders",
			args: args{
				ids: reader.AccountID{
					CompoundID: test.MyProfile.FtcID,
					FtcID:      null.StringFrom(test.MyProfile.FtcID),
					UnionID:    null.StringFrom(test.MyProfile.UnionID),
				},
				p: gorest.NewPagination(1, 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListOrders(tt.args.ids, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("ListStaff Orders: %+v", got)
		})
	}
}

func TestEnv_RetrieveOrder(t *testing.T) {

	env := Env{DB: test.DBX}

	order := test.MyProfile.Order(true)

	if err := env.CreateOrder(order); err != nil {
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
			name:    "Retrieve Order",
			args:    args{id: order.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveOrder(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Retrieved order: %+v", got)
		})
	}
}
