package registry

import (
	gorest "github.com/FTChinese/go-rest"
	oauth2 "github.com/FTChinese/superyard/internal/pkg/oauth"
	"log"
)

// CreateApp registers a new app.
func (env Env) CreateApp(app oauth2.App) error {

	_, err := env.dbs.Write.NamedExec(oauth2.StmtInsertApp, app)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) countApp() (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, oauth2.StmtCountApp)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listApps(p gorest.Pagination) ([]oauth2.App, error) {

	apps := make([]oauth2.App, 0)
	err := env.dbs.Read.Select(
		&apps,
		oauth2.StmtListApps,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return apps, nil
}

// ListApps retrieves all apps for next-api with pagination support.
func (env Env) ListApps(p gorest.Pagination) (oauth2.AppList, error) {
	countCh := make(chan int64)
	listCh := make(chan oauth2.AppList)

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

		listCh <- oauth2.AppList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return oauth2.AppList{}, listResult.Err
	}

	return oauth2.AppList{
		Total:      count,
		Pagination: p,
		Data:       listResult.Data,
		Err:        nil,
	}, nil
}

// RetrieveApp retrieves an ftc app regardless of who owns it.
// The whole team should be accessible to all apps.
func (env Env) RetrieveApp(clientID string) (oauth2.App, error) {

	var app oauth2.App
	err := env.dbs.Read.Get(&app, oauth2.StmtApp, clientID)

	if err != nil {
		return app, err
	}

	return app, nil
}

// UpdateApp allows user to update a ftc app.
func (env Env) UpdateApp(app oauth2.App) error {

	_, err := env.dbs.Read.NamedExec(oauth2.StmtUpdateApp, app)

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

	_, err = tx.Exec(oauth2.StmtRemoveApp, clientID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(oauth2.StmtRemoveAppKeys, clientID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
