package model

import (
	"database/sql"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/apn"
	"strings"
)

type APNEnv struct {
	DB *sql.DB
}

func (env APNEnv) LatestStoryList() ([]apn.Teaser, error) {
	query := storyTeaser + `
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

	teasers := make([]apn.Teaser, 0)

	for rows.Next() {
		var t apn.Teaser
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

func (env APNEnv) FindStory(id string) (apn.Teaser, error) {
	query := storyTeaser + `
	WHERE story.id = ?;`

	var t apn.Teaser
	var tag string
	err := env.DB.QueryRow(query, id).Scan(
		&t.ArticleID,
		&t.Title,
		&t.Standfirst,
		&t.Author,
		&tag,
		&t.CreatedAt,
		&t.UpdateAt,
		&t.CoverURL)

	if err != nil {
		return t, err
	}

	t.Tags = strings.Split(tag, ",")
	if !t.CoverURL.IsZero() {
		t.CoverURL = null.StringFrom(strings.Replace(
			t.CoverURL.String,
			"/upload",
			"http://i.ftimg.net",
			1))
	}
	return t, nil
}

func (env APNEnv) FindVideo(id string) (apn.Teaser, error) {
	query := `
	SELECT id AS id, 
		cheadline AS title,
		clongleadbody AS standfirst,
		CONCAT(cdescribe, ' ', cbyline) AS author,
		cc_piclink AS coverUrl
	FROM cmstmp01.video_story
	WHERE id = ? AND publish_status = 'publish'`

	var t apn.Teaser
	err := env.DB.QueryRow(query, id).Scan(
		&t.ArticleID,
		&t.Title,
		&t.Standfirst,
		&t.Author,
		&t.CoverURL)

	if err != nil {
		return t, err
	}

	return t, nil
}

func (env APNEnv) FindGallery(id string) (apn.Teaser, error) {
	query := `
	SELECT photonewsid AS id, 
		cn_title AS title,
		leadbody AS standfirst,
		tags,
		FROM_UNIXTIME(add_times) AS createdAt,
		thumb_url AS coverUrl
	FROM cmstmp01.photonews
	WHERE photonewsid = ?`

	var t apn.Teaser
	var tag string
	err := env.DB.QueryRow(query, id).Scan(
		&t.ArticleID,
		&t.Title,
		&t.Standfirst,
		&tag,
		&t.CreatedAt,
		&t.CoverURL)

	if err != nil {
		return t, err
	}

	t.Tags = strings.Split(tag, ",")
	t.UpdateAt = t.CreatedAt

	if !t.CoverURL.IsZero() {
		t.CoverURL = null.StringFrom("http://i.ftimg.net/" + t.CoverURL.String)
	}

	return t, nil
}

func (env APNEnv) FindInteractive(id string) (apn.Teaser, error) {
	query := `
	SELECT id AS id, 
		cheadline AS title,
		clongleadbody AS standfirst,
		CONCAT(cbyline_description, ' ', cauthor) AS author,
		tag,
		FROM_UNIXTIME(fileupdatetime) AS updatedAt
	FROM cmstmp01.interactive_story
	WHERE id = ?`

	var t apn.Teaser
	var tag string
	err := env.DB.QueryRow(query, id).Scan(
		&t.ArticleID,
		&t.Title,
		&t.Standfirst,
		&t.Author,
		&tag,
		&t.UpdateAt)

	if err != nil {
		return t, err
	}

	t.Tags = strings.Split(tag, ",")
	t.CreatedAt = t.UpdateAt

	return t, nil
}
