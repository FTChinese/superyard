package android

const releaseCols = `
SELECT version_name,
	version_code,
	body,
	apk_url,
	created_utc AS created_at,
	updated_utc AS updated_at
FROM file_store.android_release`

const StmtRelease = releaseCols + `
WHERE version_name = ?
LIMIT 1`

const StmtInsertRelease = `
INSERT INTO file_store.android_release
SET version_name = :version_name,
	version_code = :version_code,
	body = :body,
	apk_url = :apk_url,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const StmtUpdateRelease = `
UPDATE file_store.android_release
SET version_code = :version_code,
	body = :body,
	apk_url = :apk_url,
	updated_utc = UTC_TIMESTAMP()
WHERE version_name = :version_name
LIMIT 1`

const StmtListRelease = releaseCols + `
ORDER BY version_code DESC
LIMIT ? OFFSET ?`

const StmtReleaseExists = `
SELECT EXISTS (
	SELECT *
	FROM file_store.android_release
	WHERE version_name = ?
) AS already_exists`

const StmtDeleteRelease = `
DELETE FROM file_store.android_release
WHERE version_name = ?
LIMIT 1`
