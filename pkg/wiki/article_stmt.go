package wiki

const articleTeaserCols = `
SELECT id,
	author,
	created_utc,
	updated_utc,
	title,
	summary,
	keyword
`

const articleCols = articleTeaserCols + `,
	body
FROM file_store.wiki`

const StmtInsertArticle = `
INSERT INTO file_store.wiki
SET id = :id,
	title = :title,
	author = :author,
	summary = :summary,
	body = :body,
	keyword = :keyword,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const StmtUpdateArticle = `
UPDATE file_store.wiki
SET title = :title,
	summary = :summary,
	body = :body,
	keyword = :keyword,
	updated_utc = UTC_TIMESTAMP()
WHERE id = :id
LIMIT 1`

const StmtArticle = articleCols + `
WHERE id = ?
LIMIT 1`

const StmtListArticle = articleTeaserCols + `
FROM file_store.wiki
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`
