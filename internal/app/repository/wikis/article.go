package wikis

import (
	gorest "github.com/FTChinese/go-rest"
	wiki2 "github.com/FTChinese/superyard/internal/pkg/wiki"
)

// CreateArticle creates a new article.
func (env Env) CreateArticle(a wiki2.Article) (int64, error) {
	result, err := env.dbs.Write.NamedExec(wiki2.StmtInsertArticle, a)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (env Env) UpdateArticle(a wiki2.Article) error {
	_, err := env.dbs.Write.NamedExec(wiki2.StmtUpdateArticle, a)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadArticle(id int64) (wiki2.Article, error) {
	var a wiki2.Article
	if err := env.dbs.Read.Get(&a, wiki2.StmtArticle, id); err != nil {
		return wiki2.Article{}, err
	}

	return a, nil
}

func (env Env) ListArticles(p gorest.Pagination) ([]wiki2.Article, error) {
	var articles = make([]wiki2.Article, 0)

	err := env.dbs.Read.Select(
		&articles,
		wiki2.StmtListArticle,
		p.Limit,
		p.Offset(),
	)
	if err != nil {
		return articles, err
	}

	return articles, nil
}
