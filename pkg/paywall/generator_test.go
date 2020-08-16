package paywall

import "testing"

func Test_genProductID(t *testing.T) {
	t.Log(GenProductID())
}

func Test_genPlanID(t *testing.T) {
	t.Log(genPlanID())
}

func Test_genDiscountID(t *testing.T) {
	t.Log(genDiscountID())
}

func Test_genPromoID(t *testing.T) {
	t.Log(genPromoID())
}
