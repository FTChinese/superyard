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
FROM backyard.wiki`

const StmtInsertArticle = `
INSERT INTO backyard.wiki
SET id = :id,
	title = :title,
	author = :author,
	summary = :summary,
	body = :body,
	keyword = :keyword,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const StmtUpdateArticle = `
UPDATE backyard.wiki
SET title = :title,
	summary = :summary,
	body = :body,
	keyword = :keyword,
	updated_utc = UTC_TIMESTAMP()`

const StmtArticle = articleCols + `
WHERE id = ?
LIMIT 1`

const StmtListArticle = articleTeaserCols + `
FROM backyard.wiki
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`
