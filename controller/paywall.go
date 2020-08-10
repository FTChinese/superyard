package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type PaywallRouter struct {
}

func NewPaywallRouter(db *sqlx.DB) PaywallRouter {
	return PaywallRouter{}
}

// CreateBanner creates a single unique banner.
func (router PaywallRouter) CreateBanner(c echo.Context) error {
	return nil
}

func (router PaywallRouter) LoadBanner(c echo.Context) error {
	return nil
}

func (router PaywallRouter) UpdateBanner(c echo.Context) error {
	return nil
}

// CreatePromo creates a new promo and apply it to current banner.
func (router PaywallRouter) CreatePromo(c echo.Context) error {
	return nil
}

func (router PaywallRouter) LoadPromo(c echo.Context) error {
	return nil
}

// LoadProducts retrieves all active products presented on paywall,
// together with each product's plans, and each plan's optional
// discount.
func (router PaywallRouter) LoadProducts(c echo.Context) error {
	return nil
}
