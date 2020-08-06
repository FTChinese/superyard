package product

import "testing"

func Test_genProductID(t *testing.T) {
	t.Log(genProductID())
}

func Test_genPlanID(t *testing.T) {
	t.Log(genPlanID())
}

func Test_genDiscountID(t *testing.T) {
	t.Log(genDiscountID())
}
