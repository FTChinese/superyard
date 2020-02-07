package apps

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/superyard/models/android"
	"gitlab.com/ftchinese/superyard/models/builder"
)

type AndroidEnv struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository/apps")

func (env AndroidEnv) Exists(tag string) (bool, error) {
	var ok bool
	err := env.DB.Get(&ok, ReleaseExists, tag)

	if err != nil {
		return false, err
	}

	return ok, nil
}

func (env AndroidEnv) CreateRelease(r android.Release) error {
	_, err := env.DB.NamedExec(
		InsertRelease,
		r)

	if err != nil {
		return err
	}

	return nil
}

func (env AndroidEnv) ListReleases(p builder.Pagination) ([]android.Release, error) {
	releases := make([]android.Release, 0)

	err := env.DB.Select(
		&releases,
		listRelease,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "AndroidEnv.ListReleases").Error(err)

		return nil, err
	}

	return releases, nil
}

func (env AndroidEnv) RetrieveRelease(versionName string) (android.Release, error) {
	var r android.Release

	err := env.DB.Get(&r, selectAnRelease, versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.RetrieveRelease").Error(err)
		return r, err
	}

	return r, nil
}

func (env AndroidEnv) UpdateRelease(r android.Release) error {
	_, err := env.DB.NamedExec(
		UpdateRelease,
		r)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.UpdateRelease").Error(err)
		return err
	}

	return nil
}

func (env AndroidEnv) DeleteRelease(versionName string) error {
	_, err := env.DB.Exec(DeleteRelease, versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.DeleteRelease").Error(err)
		return err
	}

	return nil
}
