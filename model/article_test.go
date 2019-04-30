package model

import (
	"database/sql"
	"testing"
)

func TestArticleEnv_LatestCover(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Latest Cover Articles",
			fields:  fields{DB: db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := ArticleEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LatestCover()
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleEnv.LatestStoryList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Teasers: %+v", got)
		})
	}
}
