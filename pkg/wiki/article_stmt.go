package wiki

const articleBaseCols = `
SELECT id,
	title,
	author,
	summary,
	body,
	keyword,
	created_utc,
	updated_utc
`

const articleCols = articleBaseCols + `,
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
SET title = :title
	summary = :summary
	body = :body,
	keyword = :keyword,
	updated_utc = UTC_TIMESTAMP()`

const StmtArticle = articleCols + `
WHERE id = ?
LIMIT 1`

const StmtListArticle = articleCols + `
FROM backyard.wiki
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`
