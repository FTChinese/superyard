package price

import "github.com/guregu/null"

// FtcPriceParams is the form data submitted to create a price.
// A new plan is always created under a certain product.
// Therefore, the input data does not have tier field.
type FtcPriceParams struct {
	CreatedBy string `json:"createdBy"`
	Edition
	Description null.String `json:"description"`
	LiveMode    bool        `json:"liveMode"`
	Nickname    null.String `json:"nickname"`
	Price       float64     `json:"price"`
	ProductID   string      `json:"productId"`
}
