package price

import (
	"github.com/FTChinese/superyard/pkg/dt"
	"github.com/guregu/null"
)

type OfferKind string

const (
	OfferKindNull         OfferKind = ""
	OfferKindPromotion    OfferKind = "promotion"    // Apply to all uses
	OfferKindRetention    OfferKind = "retention"    // Apply only to valid user
	OfferKindWinBack      OfferKind = "win_back"     // Apply only to expired user
	OfferKindIntroductory OfferKind = "introductory" // Apply only to a new user who has not enjoyed an introductory offer
)

// DiscountParams contains fields submitted by client
// when creating a discount.
type DiscountParams struct {
	CreatedBy   string      `json:"-" db:"created_by"`
	Description null.String `json:"description" db:"discount_desc"`
	Kind        OfferKind   `json:"kind" db:"kind"`
	Percent     null.Int    `json:"percent" db:"percent"`
	dt.DateTimePeriod
	PriceOff  null.Float `json:"priceOff" db:"price_off"`
	PriceID   string     `json:"priceId" db:"price_id"`
	Recurring bool       `json:"recurring" db:"recurring"`
}
