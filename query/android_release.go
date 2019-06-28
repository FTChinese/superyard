package query

const (
	ReleaseExists = `
	SELECT EXISTS (
		SELECT *
		FROM file_store.android_release
		WHERE version_name = ?
	) AS alreadyExists;`

	InsertRelease = `
	INSERT INTO file_store.android_release
		SET version_name = ?,
			version_code = ?,
			body = ?,
			apk_url = ?,
			created_utc = UTC_TIMESTAMP(),
			updated_utc = UTC_TIMESTAMP()`

	androidRelease = `SELECT version_name,
		version_code,
		body,
		apk_url,
		created_utc,
		updated_utc
	FROM file_store.android_release`

	SingleRelease = androidRelease + `
	WHERE version_name = ?
	LIMIT 1`

	AllReleases = androidRelease + `
	ORDER BY version_code DESC
	LIMIT ? OFFSET ?`

	UpdateRelease = `
	UPDATE file_store.android_release
	SET version_code = ?,
		body = ?,
		apk_url = ?,
		updated_utc = UTC_TIMESTAMP()
	WHERE version_name = ?
	LIMIT 1`

	DeleteRelease = `
	DELETE FROM file_store.android_release
	WHERE version_name = ?
	LIMIT 1`
)
