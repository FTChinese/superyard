package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

// ListProducts retrieves a list of products with plans attached.
// The plans attached are only used for display purpose.
func (router PaywallRouter) ListProducts(c echo.Context) error {
	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).
		ListProduct(claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// CreateProduct creates a new product.
// Request body:
// - createdBy: string;
// - description?: string;
// - heading: string;
// - smallPrint?: string;
// - tier: standard | premium;
func (router PaywallRouter) CreateProduct(c echo.Context) error {

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClients.
		Select(live).
		CreateProduct(c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// LoadProduct retrieves a single product used when display
// details of a product, or editing it.
func (router PaywallRouter) LoadProduct(c echo.Context) error {
	productID := c.Param("productId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.Select(live).LoadProduct(productID, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// UpdateProduct modifies a product.
// Input
// tier: string;
// heading: string;
// description?: string;
// smallPrint?: string;
func (router PaywallRouter) UpdateProduct(c echo.Context) error {
	id := c.Param("productId")
	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClients.
		Select(live).
		UpdateProduct(id, c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// ActivateProduct puts a product on paywall.
// Request empty.
func (router PaywallRouter) ActivateProduct(c echo.Context) error {
	prodID := c.Param("productId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.Select(live).ActivateProduct(prodID, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
