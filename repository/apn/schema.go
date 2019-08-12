package apn

const (
	storyTeaser = `
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
		AND picture.id = storyToPic.picture_id`
)
