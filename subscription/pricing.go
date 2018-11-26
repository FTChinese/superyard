package subscription

import "encoding/json"

// Plan contains details of subscription plan.
type Plan struct {
	Tier  string  `json:"tier"`
	Cycle string  `json:"cycle"`
	Price float64 `json:"price"`
	ID    int
	// For wxpay, this is used as `body` parameter;
	// For alipay, this is used as `subject` parameter.
	Description string `json:"description"` // required, max 128 chars
	// For wxpay, this is used as `detail` parameter;
	// For alipay, this is used as `body` parameter.
	Message string `json:"message"`
}

// SavePricing set the pricing plans of a promotion schedule.
func (env Env) SavePricing(id int64, plans map[string]Plan) error {
	query := `
	UPDATE premium.promotion_schedule
	SET plans = ?
	WEHRE id = ?
	LIMIT 1`

	p, err := json.Marshal(plans)

	if err != nil {
		logger.WithField("location", "NewPricing").Error(err)
		return err
	}

	_, err = env.DB.Exec(query, string(p))

	if err != nil {
		logger.WithField("location", "NewPricing").Error(err)
		return err
	}

	return nil
}
