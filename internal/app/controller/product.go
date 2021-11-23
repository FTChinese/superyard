package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/products"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type ProductRouter struct {
	repo       products.Env
	apiClients subsapi.APIClients
	logger     *zap.Logger
}

func NewProductRouter(myDBs db.ReadWriteMyDBs, clients subsapi.APIClients, logger *zap.Logger) ProductRouter {
	return ProductRouter{
		repo:       products.NewEnv(myDBs),
		apiClients: clients,
		logger:     logger,
	}
}

// ListPricedProducts retrieves a list of products with plans attached.
// The plans attached are only used for display purpose.
// Deprecated
func (router ProductRouter) ListPricedProducts(c echo.Context) error {
	prods, err := router.repo.ListProducts()
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
//
// In case of 422 error, the fields in plans should
// look like:
// {
//	"field": "plans.0.price",
//	"code": "invalid"
// }
// Deprecated
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
// Deprecated
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
// Deprecated
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

	// Only products the plans activated could be set on paywall.
	// If a product has plans, but none of them is activated, then the product will have no plans
	// when presented on paywall. Thus we cannot allow it to be set on paywall.
	ok, err := router.repo.ProductHasActivePlan(prodID)
	if err != nil {
		return render.NewDBError(err)
	}

	if !ok {
		return render.NewUnprocessable(&render.ValidationError{
			Message: "This product does not have active prices set yet",
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

func (router ProductRouter) CreatePrice(c echo.Context) error {
	resp, err := router.apiClients.Live.CreatePrice(c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ProductRouter) ListPriceOfProduct(c echo.Context) error {
	productID := c.QueryParam("product_id")
	if productID == "" {
		return render.NewBadRequest("Missing query parameter product_id")
	}

	resp, err := router.apiClients.Live.ListPriceOfProduct(productID)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ProductRouter) ActivatePrice(c echo.Context) error {
	id := c.Param("planId")

	resp, err := router.apiClients.Live.ActivatePrice(id)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ProductRouter) RefreshPriceDiscounts(c echo.Context) error {
	id := c.Param("planId")

	resp, err := router.apiClients.Live.RefreshPriceDiscounts(id)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// CreateDiscount creates a discount for a plan and apply to
// that plan immediately.
// Deprecated
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

func (router ProductRouter) CreateDiscountV2(c echo.Context) error {
	resp, err := router.apiClients.Live.CreateDiscount(c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ProductRouter) RemoveDiscount(c echo.Context) error {
	id := c.Param("id")

	resp, err := router.apiClients.Live.RemoveDiscount(id)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
