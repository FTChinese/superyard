package paywall

// ExpandedProduct defines the complete data of a product.
// This is used as the products field in a paywall.
// separately and assemble them.
type ExpandedProduct struct {
	Product
	Plans []ExpandedPlan `json:"plans"`
}

// BuildPaywallProducts zips product with its plans.
func BuildPaywallProducts(prods []Product, plans []ExpandedPlan) []ExpandedProduct {
	groupedPlans := GroupPlans(plans)

	var result = make([]ExpandedProduct, 0)

	for _, prod := range prods {
		gPlans, ok := groupedPlans[prod.ID]

		if !ok {
			gPlans = []ExpandedPlan{}
		}

		result = append(result, ExpandedProduct{
			Product: prod,
			Plans:   gPlans,
		})
	}

	return result
}

type Paywall struct {
	Banner   Banner            `json:"banner"`
	Promo    Promo             `json:"promo"`
	Products []ExpandedProduct `json:"products"`
}
