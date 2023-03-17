package wikis

import (
	"testing"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/internal/pkg/wiki"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
)

func TestEnv_CreateArticle(t *testing.T) {
	env := NewEnv(db.MustNewMyDBs(false))
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		a wiki.Article
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create article",
			args: args{
				a: test.NewArticle(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.CreateArticle(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func mustCreateArticle() wiki.Article {
	env := NewEnv(db.MustNewMyDBs(false))
	article := test.NewArticle()
	id, err := env.CreateArticle(article)
	if err != nil {
		panic(err)
	}

	article.ID = id

	return article
}

func TestEnv_LoadArticle(t *testing.T) {
	article := mustCreateArticle()

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		//want    wiki.Article
		wantErr bool
	}{
		{
			name: "Load an article",
			args: args{
				id: article.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadArticle(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadArticle() got = %v, want %v", got, tt.want)
			//}
			t.Logf("%v", got)
		})
	}
}

func TestEnv_ListArticles(t *testing.T) {
	mustCreateArticle()

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name string
		args args
		//want    []wiki.ArticleTeaser
		wantErr bool
	}{
		{
			name: "List articles",
			args: args{
				p: gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListArticles(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ListArticles() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%v", got)
		})
	}
}
