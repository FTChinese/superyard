package apps

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/android"
	"log"
)

// CreateRelease insert a new row of android release.
func (env Env) CreateRelease(r android.Release) error {
	_, err := env.DB.NamedExec(
		android.StmtInsertRelease,
		r)

	if err != nil {
		return err
	}

	return nil
}

// RetrieveRelease retrieves a row of release.
func (env Env) RetrieveRelease(versionName string) (android.Release, error) {
	var r android.Release

	err := env.DB.Get(&r, android.StmtRelease, versionName)

	if err != nil {
		return r, err
	}

	return r, nil
}

// UpdateRelease updates a release.
func (env Env) UpdateRelease(input android.ReleaseInput) error {
	_, err := env.DB.NamedExec(
		android.StmtUpdateRelease,
		input)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) countRelease() (int64, error) {
	var count int64
	err := env.DB.Get(&count, android.StmtCountRelease)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listReleases(p gorest.Pagination) ([]android.Release, error) {
	releases := make([]android.Release, 0)

	err := env.DB.Select(
		&releases,
		android.StmtListRelease,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return releases, nil
}

// ListRelease lists all releases.
func (env Env) ListReleases(p gorest.Pagination) (android.ReleaseList, error) {
	countCh := make(chan int64)
	listCh := make(chan android.ReleaseList)

	go func() {
		defer close(countCh)
		n, err := env.countRelease()
		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listReleases(p)
		listCh <- android.ReleaseList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return android.ReleaseList{}, listResult.Err
	}

	return android.ReleaseList{
		Total:      count,
		Pagination: p,
		Data:       listResult.Data,
		Err:        nil,
	}, nil
}

// Exists checks whether a release already exists.
func (env Env) Exists(tag string) (bool, error) {
	var ok bool
	err := env.DB.Get(&ok, android.StmtReleaseExists, tag)

	if err != nil {
		return false, err
	}

	return ok, nil
}

// Delete a release removes a release.
func (env Env) DeleteRelease(versionName string) error {
	_, err := env.DB.Exec(android.StmtDeleteRelease, versionName)

	if err != nil {
		return err
	}

	return nil
}
