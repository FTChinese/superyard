package apps

import (
	"gitlab.com/ftchinese/superyard/models/android"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/stmt"
)

const stmtInsertRelease = `
INSERT INTO file_store.android_release
SET version_name = :version_name,
	version_code = :version_code,
	body = :body,
	apk_url = :apk_url,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

func (env AndroidEnv) CreateRelease(r android.Release) error {
	_, err := env.DB.NamedExec(
		stmtInsertRelease,
		r)

	if err != nil {
		return err
	}

	return nil
}

const selectAnRelease = stmt.AndroidRelease + `
WHERE version_name = ?
LIMIT 1`

func (env AndroidEnv) RetrieveRelease(versionName string) (android.Release, error) {
	var r android.Release

	err := env.DB.Get(&r, selectAnRelease, versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.RetrieveRelease").Error(err)
		return r, err
	}

	return r, nil
}

const UpdateRelease = `
UPDATE file_store.android_release
SET version_code = :version_code,
	body = :body,
	apk_url = :apk_url,
	updated_utc = UTC_TIMESTAMP()
WHERE version_name = :version_name
LIMIT 1`

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

const listRelease = stmt.AndroidRelease + `
ORDER BY version_code DESC
LIMIT ? OFFSET ?`

func (env AndroidEnv) ListReleases(p util.Pagination) ([]android.Release, error) {
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

const stmtReleaseExists = `
SELECT EXISTS (
	SELECT *
	FROM file_store.android_release
	WHERE version_name = ?
) AS already_exists;`

func (env AndroidEnv) Exists(tag string) (bool, error) {
	var ok bool
	err := env.DB.Get(&ok, stmtReleaseExists, tag)

	if err != nil {
		return false, err
	}

	return ok, nil
}

const DeleteRelease = `
DELETE FROM file_store.android_release
WHERE version_name = ?
LIMIT 1`

func (env AndroidEnv) DeleteRelease(versionName string) error {
	_, err := env.DB.Exec(DeleteRelease, versionName)

	if err != nil {
		logger.WithField("trace", "AndroidEnv.DeleteRelease").Error(err)
		return err
	}

	return nil
}
