package stripeapi

import "github.com/stripe/stripe-go/v72"

func (c Client) ListPrices() (*stripe.PriceList, error) {
	return c.sc.Prices.List(&stripe.PriceListParams{
		Active: stripe.Bool(true),
	}).PriceList(), nil
}

func (c Client) RetrievePrice(id string) (*stripe.Price, error) {
	return c.sc.Prices.Get(id, nil)
}
