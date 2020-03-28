package stmt

const AndroidRelease = `
SELECT version_name,
	version_code,
	body,
	apk_url,
	created_utc AS created_at,
	updated_utc AS updated_at
FROM file_store.android_release`
