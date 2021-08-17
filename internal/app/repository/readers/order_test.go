package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/test"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_ListOrders(t *testing.T) {
	p := test.NewPersona()

	repo := test.NewRepo()

	repo.MustCreateOrder(p.Order(true))
	repo.MustCreateOrder(p.SetAccountKind(enum.AccountKindWx).Order(false))
	repo.MustCreateOrder(p.SetAccountKind(enum.AccountKindLinked).Order(true))

	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ids reader.IDs
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
			args: args{
				ids: reader.IDs{
					FtcID:   null.NewString(p.FtcID, p.FtcID != ""),
					UnionID: null.NewString(p.UnionID, p.UnionID != ""),
				},
				p: gorest.NewPagination(1, 10),
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.ListOrders(tt.args.ids, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, len(got.Data))

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_RetrieveOrder(t *testing.T) {
	p := test.NewPersona()
	order := p.Order(true)

	repo := test.NewRepo()
	repo.MustCreateOrder(order)

	env := NewEnv(db.MustNewMyDBs(false), zaptest.NewLogger(t))

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Loads an order",
			args:    args{id: order.ID},
			want:    order.ID,
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
			assert.Equal(t, tt.want, got.ID)
		})
	}
}
