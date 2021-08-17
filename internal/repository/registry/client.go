package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/oauth"
	"log"
)

// CreateApp registers a new app.
func (env Env) CreateApp(app oauth.App) error {

	_, err := env.dbs.Write.NamedExec(oauth.StmtInsertApp, app)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) countApp() (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, oauth.StmtCountApp)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listApps(p gorest.Pagination) ([]oauth.App, error) {

	apps := make([]oauth.App, 0)
	err := env.dbs.Read.Select(
		&apps,
		oauth.StmtListApps,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return apps, nil
}

// ListApps retrieves all apps for next-api with pagination support.
func (env Env) ListApps(p gorest.Pagination) (oauth.AppList, error) {
	countCh := make(chan int64)
	listCh := make(chan oauth.AppList)

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

		listCh <- oauth.AppList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return oauth.AppList{}, listResult.Err
	}

	return oauth.AppList{
		Total:      count,
		Pagination: p,
		Data:       listResult.Data,
		Err:        nil,
	}, nil
}

// RetrieveApp retrieves an ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env Env) RetrieveApp(clientID string) (oauth.App, error) {

	var app oauth.App
	err := env.dbs.Read.Get(&app, oauth.StmtApp, clientID)

	if err != nil {
		return app, err
	}

	return app, nil
}

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(app oauth.App) error {

	_, err := env.dbs.Read.NamedExec(oauth.StmtUpdateApp, app)

	if err != nil {
		return err
	}

	return nil
}

// RemoveApp deactivate an ftc app.
// All access tokens belonging to this app should be deactivated.
func (env Env) RemoveApp(clientID string) error {
	tx, err := env.dbs.Read.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(oauth.StmtRemoveApp, clientID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(oauth.StmtRemoveAppKeys, clientID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
