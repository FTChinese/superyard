package apn

import (
	"database/sql"
	"gitlab.com/ftchinese/backyard-api/model"
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
			fields:  fields{DB: model.db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := ArticleEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LatestStoryList()
			if (err != nil) != tt.wantErr {
				t.Errorf("APNEnv.LatestStoryList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Teasers: %+v", got)
		})
	}
}

func TestArticleEnv_FindStory(t *testing.T) {
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
			name:    "Story Teaser",
			fields:  fields{DB: model.db},
			args:    args{id: "001076300"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := ArticleEnv{
				DB: tt.fields.DB,
			}
			got, err := env.FindStory(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("APNEnv.FindStory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Story teaser: %+v", got)
		})
	}
}

func TestArticleEnv_FindVideo(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    article.Teaser
		wantErr bool
	}{
		{
			name:    "Video Teaser",
			fields:  fields{DB: model.db},
			args:    args{id: "2586"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := ArticleEnv{
				DB: tt.fields.DB,
			}
			got, err := env.FindVideo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("APNEnv.FindVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("APNEnv.FindVideo() = %v, want %v", got, tt.want)
			//}

			t.Logf("Video teaser: %+v", got)
		})
	}
}

func TestArticleEnv_FindGallery(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    article.Teaser
		wantErr bool
	}{
		{
			name:    "Gallery Teaser",
			fields:  fields{DB: model.db},
			args:    args{id: "1050"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := ArticleEnv{
				DB: tt.fields.DB,
			}
			got, err := env.FindGallery(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("APNEnv.FindGallery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("APNEnv.FindGallery() = %v, want %v", got, tt.want)
			//}

			t.Logf("Gallery teaser: %+v", got)
		})
	}
}

func TestArticleEnv_FindInteractive(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    article.Teaser
		wantErr bool
	}{
		{
			name:    "Interactive Teaser",
			fields:  fields{DB: model.db},
			args:    args{id: "5065"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := ArticleEnv{
				DB: tt.fields.DB,
			}
			got, err := env.FindInteractive(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("APNEnv.FindInteractive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("APNEnv.FindInteractive() = %v, want %v", got, tt.want)
			//}

			t.Logf("Interactive teaser: %+v", got)
		})
	}
}
