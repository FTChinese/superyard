package wikis

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/wiki"
)

// CreateArticle creates a new article.
func (env Env) CreateArticle(a wiki.Article) (int64, error) {
	result, err := env.db.NamedExec(wiki.StmtInsertArticle, a)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (env Env) UpdateArticle(a wiki.Article) error {
	_, err := env.db.NamedExec(wiki.StmtUpdateArticle, a)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadArticle(id int64) (wiki.Article, error) {
	var a wiki.Article
	if err := env.db.Get(&a, wiki.StmtArticle, id); err != nil {
		return wiki.Article{}, err
	}

	return a, nil
}

func (env Env) ListArticles(p gorest.Pagination) ([]wiki.ArticleTeaser, error) {
	var articles = make([]wiki.ArticleTeaser, 0)

	err := env.db.Select(
		&articles,
		wiki.StmtListArticle,
		p.Limit,
		p.Offset(),
	)
	if err != nil {
		return articles, err
	}

	return articles, nil
}
