package wikis

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/internal/pkg/wiki"
)

// CreateArticle creates a new article.
func (env Env) CreateArticle(a wiki.Article) (wiki.Article, error) {
	result := env.gormDBs.Write.Create(&a)

	if result.Error != nil {
		return a, result.Error
	}

	return a, nil
}

func (env Env) UpdateArticle(a wiki.Article) error {
	err := env.gormDBs.Write.Save(&a).Error

	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadArticle(id int64) (wiki.Article, error) {
	var a wiki.Article
	err := env.gormDBs.Read.Where("id = ?", id).Take(&a).Error

	if err != nil {
		return wiki.Article{}, err
	}

	return a, nil
}

func (env Env) ListArticles(p gorest.Pagination) ([]wiki.Article, error) {
	var articles = make([]wiki.Article, 0)

	err := env.gormDBs.Read.
		Order("created_utc DESC").
		Limit(int(p.Limit)).
		Offset(int(p.Offset())).
		Find(&articles).Error

	if err != nil {
		return articles, err
	}

	return articles, nil
}
