package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type ProductRouter struct {
}

func NewProductRouter(db *sqlx.DB) ProductRouter {
	return ProductRouter{}
}

func (router ProductRouter) CreateProduct(c echo.Context) error {
	return nil
}

// ListProducts retrieves a list of products with plans attached.
// The plans attached are only used for display purpose.
func (router ProductRouter) ListProducts(c echo.Context) error {
	return nil
}

// LoadProduct retrieves a single product used when display
// details of a product, or editing it.
func (router ProductRouter) LoadProduct(c echo.Context) error {
	return nil
}

func (router ProductRouter) UpdateProduct(c echo.Context) error {
	return nil
}

// CreatePlan creates a new plan for a product.
func (router ProductRouter) CreatePlan(c echo.Context) error {
	return nil
}

// ActivatePlan puts a plan as the default under a product.
// This will make the plan visible on paywall.
func (router ProductRouter) ActivatePlan(c echo.Context) error {
	return nil
}

// ListPlans retrieves all plans of a product.
func (router ProductRouter) ListPlans(c echo.Context) error {
	return nil
}

// CreateDiscount creates a discount for a plan and apply to
// that plan immediately.
func (router ProductRouter) CreateDiscount(c echo.Context) error {
	return nil
}
