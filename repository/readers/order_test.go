package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_ListOrders(t *testing.T) {
	p := test.NewPersona()

	repo := test.NewRepo()

	repo.MustCreateOrder(p.Order(true))
	repo.MustCreateOrder(p.SetAccountKind(reader.AccountKindWx).Order(false))
	repo.MustCreateOrder(p.SetAccountKind(reader.AccountKindLinked).Order(true))

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ids subs.CompoundIDs
		p   gorest.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "List orders",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				ids: subs.CompoundIDs{
					FtcID:   p.FtcID,
					UnionID: p.UnionID,
				},
				p: gorest.NewPagination(1, 10),
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.ListOrders(tt.args.ids, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, len(got))
		})
	}
}

func TestEnv_RetrieveOrder(t *testing.T) {
	p := test.NewPersona()
	order := p.Order(true)

	repo := test.NewRepo()
	repo.MustCreateOrder(order)

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Loads an order",
			fields:  fields{DB: test.DBX},
			args:    args{id: order.ID},
			want:    order.ID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.RetrieveOrder(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got.ID)
		})
	}
}

func TestEnv_ConfirmOrder(t *testing.T) {
	p := test.NewPersona()
	order := p.Order(false)

	repo := test.NewRepo()
	repo.MustCreateOrder(order)

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    subs.ConfirmationResult
		wantErr bool
	}{
		{
			name: "Confirm an order",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				id: order.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.ConfirmOrder(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfirmOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Order %s Confirmed at: %s", got.Order.ID, got.Order.ConfirmedAt)
			t.Logf("Latest membership: %+v", got.Membership)
		})
	}
}
