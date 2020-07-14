package apps

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/pkg/android"
)

func (env AndroidEnv) CreateRelease(r android.Release) error {
	_, err := env.DB.NamedExec(
		android.StmtInsertRelease,
		r)

	if err != nil {
		return err
	}

	return nil
}

func (env AndroidEnv) RetrieveRelease(versionName string) (android.Release, error) {
	var r android.Release

	err := env.DB.Get(&r, android.StmtRelease, versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.RetrieveRelease").Error(err)
		return r, err
	}

	return r, nil
}

func (env AndroidEnv) UpdateRelease(r android.Release) error {
	_, err := env.DB.NamedExec(
		android.StmtUpdateRelease,
		r)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.UpdateRelease").Error(err)
		return err
	}

	return nil
}

func (env AndroidEnv) ListReleases(p gorest.Pagination) ([]android.Release, error) {
	releases := make([]android.Release, 0)

	err := env.DB.Select(
		&releases,
		android.StmtListRelease,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "AndroidEnv.ListReleases").Error(err)

		return nil, err
	}

	return releases, nil
}

func (env AndroidEnv) Exists(tag string) (bool, error) {
	var ok bool
	err := env.DB.Get(&ok, android.StmtReleaseExists, tag)

	if err != nil {
		return false, err
	}

	return ok, nil
}

func (env AndroidEnv) DeleteRelease(versionName string) error {
	_, err := env.DB.Exec(android.StmtDeleteRelease, versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.DeleteRelease").Error(err)
		return err
	}

	return nil
}
