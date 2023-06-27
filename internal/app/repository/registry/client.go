package registry

import (
	"log"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg"
	"gorm.io/gorm"
)

// CreateApp registers a new app.
func (env Env) CreateApp(app oauth.App) (oauth.App, error) {

	err := env.gormDBs.Write.Create(&app).Error

	if err != nil {
		return oauth.App{}, err
	}

	return app, nil
}

func (env Env) countApp() (int64, error) {
	var count int64

	err := env.gormDBs.Read.
		Model(&oauth.App{}).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listApps(p gorest.Pagination) ([]oauth.App, error) {

	apps := make([]oauth.App, 0)

	err := env.gormDBs.Read.
		Limit(int(p.Limit)).
		Offset(int(p.Offset())).
		Find(&apps).
		Error

	if err != nil {
		return nil, err
	}

	return apps, nil
}

// ListApps retrieves all apps for next-api with pagination support.
func (env Env) ListApps(p gorest.Pagination) (pkg.PagedList[oauth.App], error) {
	countCh := make(chan int64)
	listCh := make(chan pkg.AsyncResult[[]oauth.App])

	go func() {
		defer close(countCh)
		n, err := env.countApp()
		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listApps(p)

		listCh <- pkg.AsyncResult[[]oauth.App]{
			Value: list,
			Err:   err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return pkg.PagedList[oauth.App]{}, listResult.Err
	}

	return pkg.PagedList[oauth.App]{
		Total:      count,
		Pagination: p,
		Data:       listResult.Value,
	}, nil
}

// RetrieveApp retrieves an ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env Env) RetrieveApp(clientID string) (oauth.App, error) {

	var app oauth.App

	err := env.gormDBs.Read.
		Where("client_id = UNHEX(?)", clientID).
		Find(&app).
		Error

	if err != nil {
		return app, err
	}

	return app, nil
}

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(app oauth.App) error {

	err := env.gormDBs.Write.
		Where("is_active = ?", true).
		Save(&app).
		Error

	if err != nil {
		return err
	}

	return nil
}

const stmtRemoveAppKeys = `
UPDATE oauth.access
	SET is_active = 0
WHERE client_id = UNHEX(?)
	AND usage_type = 'app'`

// RemoveApp deactivate an ftc app.
// All access tokens belonging to this app should be deactivated.
func (env Env) RemoveApp(app oauth.App) error {

	return env.gormDBs.Write.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("is_active = ?", true).
			Save(&app).
			Error

		if err != nil {
			return err
		}

		err = tx.Raw(stmtRemoveAppKeys, app.ClientID).Error

		if err != nil {
			return err
		}

		return nil
	})
}
