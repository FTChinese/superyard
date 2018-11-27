package subscription

import "testing"

func TestSavePricing(t *testing.T) {
	err := devEnv.SavePricing(1, mockPricing)

	if err != nil {
		t.Error(err)
	}
}
