package android

const colUpsert = `
version_name = :version_name,
version_code = :version_code,
body = :body,
apk_url = :apk_url
`

const StmtInsertRelease = `
INSERT INTO file_store.android_release
SET ` + colUpsert + `,
	created_utc = :created_at`

const StmtUpdateRelease = `
UPDATE file_store.android_release
SET ` + colUpsert + `,
	updated_utc = :updated_at
WHERE version_name = :version_name
LIMIT 1`

const colSelect = `
SELECT id,
	version_name,
	version_code,
	body,
	apk_url,
	created_utc AS created_at,
	updated_utc AS updated_at
FROM file_store.android_release
`

const StmtSelectRelease = colSelect + `
WHERE version_name = ?
LIMIT 1
`

const StmtListRelease = colSelect + `
ORDER BY version_code DESC
LIMIT ? OFFSET ?
`

const StmtCountRelease = `
SELECT COUNT(*)
FROM file_store.android_release`

const StmtReleaseExists = `
SELECT EXISTS (
	SELECT *
	FROM file_store.android_release
	WHERE version_name = ?
) AS already_exists`

const StmtDeleteRelease = `
DELETE FROM file_store.android_release
WHERE version_name = ?
LIMIT 1
`
