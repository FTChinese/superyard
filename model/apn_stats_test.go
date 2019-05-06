package model

import (
	"database/sql"
	"testing"

	gorest "github.com/FTChinese/go-rest"
)

func TestAPNEnv_ListMessage(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    []article.MessageTeaser
		wantErr bool
	}{
		{
			name:    "List APN Messages",
			fields:  fields{DB: db},
			args:    args{p: gorest.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := APNEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListMessage(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("APNEnv.ListMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("APNEnv.ListMessage() = %v, want %v", got, tt.want)
			//}

			t.Logf("%+v", got)
		})
	}
}
