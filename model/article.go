package model

import (
	"database/sql"
	"gitlab.com/ftchinese/backyard-api/article"
	"strings"
)

type ArticleEnv struct {
	DB *sql.DB
}

func (env ArticleEnv) LatestCover() ([]article.Teaser, error) {
	query := `
	SELECT story.id,
		story.cheadline AS title,
		story.clongleadbody AS standfirst,
		story.cauthor AS author,
		story.tag,
		FROM_UNIXTIME(story.fileupdatetime) AS createdAt,
		FROM_UNIXTIME(story.last_publish_time) AS updatedAt,
		picture.piclink AS coverUrl
	FROM cmstmp01.story AS story
		LEFT JOIN (
			cmstmp01.story_pic AS storyToPic
			INNER JOIN cmstmp01.picture AS picture
		)
		ON story.id = storyToPic.storyid 
		AND picture.id = storyToPic.picture_id
	WHERE story.pubdate = (SELECT pubdate
			FROM cmstmp01.story
			WHERE publish_status = 'publish'
			ORDER BY pubdate DESC
			LIMIT 1)
		AND story.publish_status = 'publish'
	ORDER BY story.last_publish_time DESC`

	rows, err := env.DB.Query(query)

	if err != nil {
		logger.WithField("trace", "LatestStoryList").Error(err)
		return nil, err
	}

	defer rows.Close()

	teasers := make([]article.Teaser, 0)

	for rows.Next() {
		var t article.Teaser
		var tag string
		err := rows.Scan(
			&t.ArticleID,
			&t.Title,
			&t.Standfirst,
			&t.Author,
			&tag,
			&t.CreatedAt,
			&t.UpdateAt,
			&t.CoverURL)

		if err != nil {
			logger.WithField("trace", "LatestStoryList").Error(err)
			continue
		}

		t.Tags = strings.Split(tag, ",")

		teasers = append(teasers, t)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "LatestStoryList").Error(err)
		return teasers, err
	}

	return teasers, nil
}
