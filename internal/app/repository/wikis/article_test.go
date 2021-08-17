package wikis

import (
	gorest "github.com/FTChinese/go-rest"
	wiki2 "github.com/FTChinese/superyard/internal/pkg/wiki"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_CreateArticle(t *testing.T) {
	env := NewEnv(db.MustNewMyDBs(false))
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		a wiki2.Article
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

			assert.NotZero(t, got, "last insert id should not be 0")
		})
	}
}

func mustCreateArticle() wiki2.Article {
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

			assert.Equal(t, got.ID, article.ID, "should got the same article")
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

			assert.GreaterOrEqual(t, len(got), 1, "should retrieve at least one article")
		})
	}
}