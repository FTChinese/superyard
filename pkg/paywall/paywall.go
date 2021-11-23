package paywall

// ExpandedProduct defines the complete data of a product.
// This is used as the products field in a paywall.
// separately and assemble them.
type ExpandedProduct struct {
	Product
	Plans []ExpandedPlan `json:"plans"`
}

type Paywall struct {
	Banner   Banner            `json:"banner"`
	Promo    Promo             `json:"promo"`
	Products []ExpandedProduct `json:"products"`
}
