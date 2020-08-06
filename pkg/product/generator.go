package product

import "github.com/FTChinese/go-rest/rand"

func genProductID() string {
	return "prod_" + rand.String(12)
}

func genPlanID() string {
	return "plan_" + rand.String(12)
}

func genDiscountID() string {
	return "dsc_" + rand.String(12)
}
