package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/repository/products"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProductRouter struct {
	repo products.Env
}

func NewProductRouter(db *sqlx.DB) ProductRouter {
	return ProductRouter{
		repo: products.NewEnv(db),
	}
}

// ListPricedProducts retrieves a list of products with plans attached.
// The plans attached are only used for display purpose.
func (router ProductRouter) ListPricedProducts(c echo.Context) error {
	prods, err := router.repo.ListPricedProducts()
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, prods)
}

// CreateProduct handle request to create a new product, with optional prices.
// Input:
// tier: string;
// heading: string;
// description?: string;
// smallPrint?: string;
// plans?: [
// 	{
//		price: number;
//		cycle: string;
//		description?: string;
//	}
// ]
func (router ProductRouter) CreateProduct(c echo.Context) error {
	claims := getPassportClaims(c)
	var input paywall.PricedProductInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	product := paywall.NewPricedProduct(input, claims.Username)

	if ve := product.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	if err := router.repo.CreatePricedProduct(product); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, product)
}

// LoadProduct retrieves a single product used when display
// details of a product, or editing it.
func (router ProductRouter) LoadProduct(c echo.Context) error {
	productID := c.Param("productId")

	prod, err := router.repo.LoadProduct(productID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, prod)
}

// UpdateProduct modifies a product.
// Input
// tier: string;
// heading: string;
// description?: string;
// smallPrint?: string;
func (router ProductRouter) UpdateProduct(c echo.Context) error {
	id := c.Param("productId")

	var input paywall.ProductInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	// Use id to retrieve product.
	prod, err := router.repo.LoadProduct(id)
	if err != nil {
		return render.NewDBError(err)
	}

	// Update product
	updated := prod.Update(input)

	// Save modifications
	err = router.repo.UpdateProduct(updated)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, updated)
}

// ActivateProduct puts a product on paywall.
// Request empty.
// Response paywall.Product
func (router ProductRouter) ActivateProduct(c echo.Context) error {
	prodID := c.Param("productId")

	ok, err := router.repo.ProductHasActivePlan(prodID)
	if err != nil {
		return render.NewDBError(err)
	}

	if !ok {
		return render.NewUnprocessable(&render.ValidationError{
			Message: "This product does not have prices yet",
			Field:   "plans",
			Code:    render.CodeMissing,
		})
	}

	product, err := router.repo.LoadProduct(prodID)
	if err != nil {
		return render.NewDBError(err)
	}

	err = router.repo.ActivateProduct(product)
	if err != nil {
		return render.NewDBError(err)
	}

	product.IsActive = true
	return c.JSON(http.StatusOK, product)
}

// CreatePlan creates a new plan for a product.
// Input:
// productId: string;
// price: number;
// tier: string;
// cycle: string;
// description?: string;
func (router ProductRouter) CreatePlan(c echo.Context) error {
	claims := getPassportClaims(c)

	var input paywall.PlanInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	plan := paywall.NewPlan(input, claims.Username)

	if err := router.repo.CreatePlan(plan); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, plan)
}

// ActivatePlan puts a plan as the default under a product.
// This will make the plan visible on paywall.
func (router ProductRouter) ActivatePlan(c echo.Context) error {
	id := c.Param("planId")

	// Retrieve Plan by the id
	plan, err := router.repo.LoadPlan(id)
	if err != nil {
		return render.NewDBError(err)
	}

	// Put it into active plan table.
	err = router.repo.ActivatePlan(plan)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListPlansOfProduct retrieves all plans of a product.
// Each plan is a ExpandedPlan instance.
func (router ProductRouter) ListPlansOfProduct(c echo.Context) error {
	productID := c.QueryParam("product_id")
	if productID == "" {
		return render.NewBadRequest("Missing query parameter product_id")
	}

	plans, err := router.repo.ListPlansOfProduct(productID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, plans)
}

// CreateDiscount creates a discount for a plan and apply to
// that plan immediately.
func (router ProductRouter) CreateDiscount(c echo.Context) error {
	claims := getPassportClaims(c)
	planID := c.Param("planId")

	var input paywall.DiscountInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	schema := paywall.NewDiscountSchema(input, planID, claims.Username)

	if err := router.repo.CreateDiscount(schema); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, schema.Discount)
}

func (router ProductRouter) DropDiscount(c echo.Context) error {
	planID := c.Param("planId")

	// Retrieve Plan by the id
	plan, err := router.repo.LoadPlan(planID)
	if err != nil {
		return render.NewDBError(err)
	}

	err = router.repo.DropDiscount(plan)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
