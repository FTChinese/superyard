package model

import (
	"database/sql"
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/android"
	"gitlab.com/ftchinese/backyard-api/query"
)

type AndroidEnv struct {
	DB *sql.DB
}

func (env AndroidEnv) Exists(tag string) (bool, error) {
	var ok bool
	err := env.DB.QueryRow(query.ReleaseExists, tag).Scan(&ok)

	if err != nil {
		return false, err
	}

	return ok, nil
}

func (env AndroidEnv) CreateRelease(r android.Release) error {
	_, err := env.DB.Exec(query.InsertRelease,
		r.VersionName,
		r.VersionCode,
		r.Body,
		r.BinaryURL)

	if err != nil {
		return err
	}

	return nil
}

func (env AndroidEnv) ListReleases(p gorest.Pagination) ([]android.Release, error) {
	rows, err := env.DB.Query(
		query.AllReleases,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "AndroidEnv.ListReleases").Error(err)

		return nil, err
	}
	defer rows.Close()

	releases := make([]android.Release, 0)

	for rows.Next() {
		var r android.Release

		err := rows.Scan(
			&r.VersionName,
			&r.VersionCode,
			&r.Body,
			&r.BinaryURL,
			&r.CreatedAt,
			&r.UpdatedAt)

		if err != nil {
			logger.WithField("trace", "AndroidEnv.ListReleases").Error(err)
			return nil, err
		}

		releases = append(releases, r)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "AndroidEnv.ListReleases").Error(err)

		return releases, err
	}

	return releases, nil
}

func (env AndroidEnv) SingleRelease(versionName string) (android.Release, error) {
	var r android.Release

	err := env.DB.QueryRow(query.SingleRelease, versionName).Scan(
		&r.VersionName,
		&r.VersionCode,
		&r.Body,
		&r.BinaryURL,
		&r.CreatedAt,
		&r.UpdatedAt)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.SingleRelease").Error(err)
		return r, err
	}

	return r, nil
}

func (env AndroidEnv) UpdateRelease(r android.Release, versionName string) error {
	_, err := env.DB.Exec(query.UpdateRelease,
		r.VersionCode,
		r.Body,
		r.BinaryURL,
		versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.UpdateRelease").Error(err)
		return err
	}

	return nil
}

func (env AndroidEnv) DeleteRelease(versionName string) error {
	_, err := env.DB.Exec(query.DeleteRelease, versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.DeleteRelease").Error(err)
		return err
	}

	return nil
}
