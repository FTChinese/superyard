package subscription

import "encoding/json"

// Banner is the content used on promotion banner
type Banner struct {
	Heading    string   `json:"heading"`
	SubHeading string   `json:"subHeading"`
	Content    []string `json:"content"`
}

func (env Env) NewBanner(id int64, banner Banner) error {
	query := `
	UPDATE premium.promotion_schedule
	SET banner = ?
	WHERE id = ?
	LIMIT 1`

	b, err := json.Marshal(banner)

	if err != nil {
		logger.WithField("location", "NewBanner").Error(err)

		return err
	}

	_, err = env.DB.Exec(query, string(b))

	if err != nil {
		logger.WithField("location", "NewBanner").Error(err)
		return err
	}

	return nil
}
