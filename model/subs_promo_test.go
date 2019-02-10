package model

import (
	"database/sql"
	"log"
	"testing"

	"gitlab.com/ftchinese/backyard-api/subs"
	"gitlab.com/ftchinese/backyard-api/util"
)

func TestPromoEnv_NewSchedule(t *testing.T) {
	mStaff := newMockStaff()
	mStaff.createAccount()

	sch := mStaff.schedule()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		s       subs.Schedule
		creator string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "New Schedule",
			fields: fields{DB: db},
			args: args{
				s:       sch,
				creator: mStaff.userName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := PromoEnv{
				DB: tt.fields.DB,
			}
			got, err := env.NewSchedule(tt.args.s, tt.args.creator)
			if (err != nil) != tt.wantErr {
				t.Errorf("PromoEnv.NewSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}

func TestPromoEnv_SavePlans(t *testing.T) {
	mStaff := newMockStaff()
	mStaff.createAccount()

	id := mStaff.createSchedule()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id    int64
		plans subs.Pricing
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Save Plans",
			fields:  fields{DB: db},
			args:    args{id: id, plans: defaultPlans},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := PromoEnv{
				DB: tt.fields.DB,
			}
			if err := env.SavePlans(tt.args.id, tt.args.plans); (err != nil) != tt.wantErr {
				t.Errorf("PromoEnv.SavePlans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPromoEnv_SaveBanner(t *testing.T) {
	mStaff := newMockStaff()
	mStaff.createAccount()

	id := mStaff.createSchedule()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id     int64
		banner subs.Banner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Save Banner",
			fields:  fields{DB: db},
			args:    args{id: id, banner: mStaff.banner()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := PromoEnv{
				DB: tt.fields.DB,
			}
			if err := env.SaveBanner(tt.args.id, tt.args.banner); (err != nil) != tt.wantErr {
				t.Errorf("PromoEnv.SaveBanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPromoEnv_ListPromos(t *testing.T) {
	mStaff := newMockStaff()
	mStaff.createAccount()

	mStaff.createPromo()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		p util.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "List Promos",
			fields:  fields{DB: db},
			args:    args{p: util.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := PromoEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListPromos(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("PromoEnv.ListPromos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestPromoEnv_LoadPromo(t *testing.T) {
	mStaff := newMockStaff()
	mStaff.createAccount()

	id := mStaff.createPromo()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Load Promo",
			fields: fields{DB: db},
			args: args{id: id},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := PromoEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LoadPromo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PromoEnv.LoadPromo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}
