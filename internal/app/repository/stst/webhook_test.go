package stst

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/test"
	"testing"
)

func TestEnv_countAliUnconfirmed(t *testing.T) {

	env := NewEnv(test.DBX)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Count ali unconfirmed",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.countAliUnconfirmed()
			if (err != nil) != tt.wantErr {
				t.Errorf("countAliUnconfirmed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%d", got)
		})
	}
}

func TestEnv_listAliUnconfirmed(t *testing.T) {

	env := NewEnv(test.DBX)

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "List ali unconfirmed",
			args:    args{p: gorest.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.listAliUnconfirmed(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("listAliUnconfirmed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_AliUnconfirmed(t *testing.T) {
	env := NewEnv(test.DBX)

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "List and count ali unconfirmed",
			args: args{
				p: gorest.NewPagination(1, 20),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.AliUnconfirmed(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("AliUnconfirmed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_countWxUnconfirmed(t *testing.T) {
	env := NewEnv(test.DBX)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Count wx unconfirmed",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.countWxUnconfirmed()
			if (err != nil) != tt.wantErr {
				t.Errorf("countWxUnconfirmed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%d", got)
		})
	}
}

func TestEnv_listWxUnconfirmed(t *testing.T) {
	env := NewEnv(test.DBX)

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "List wx unconfirmed",
			args:    args{p: gorest.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.listWxUnconfirmed(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("listWxUnconfirmed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", got)
		})
	}
}

func TestEnv_WxUnconfirmed(t *testing.T) {

	env := NewEnv(test.DBX)

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Count and list wx unconfirmed",
			args: args{
				p: gorest.NewPagination(1, 20),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.WxUnconfirmed(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("WxUnconfirmed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}
