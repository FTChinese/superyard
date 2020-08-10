package paywall

import (
	"encoding/json"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/validator"
	"strings"
)

// BaseProductInput defines the request data to create a new
// product.
type BaseProductInput struct {
	Tier        enum.Tier   `json:"tier" db:"tier"`
	Heading     string      `json:"heading" db:"heading"`
	Description null.String `json:"description" db:"description"`
	SmallPrint  null.String `json:"smallPrint" db:"small_print"`
}

func (p *BaseProductInput) Validate() *render.ValidationError {
	p.Heading = strings.TrimSpace(p.Heading)
	p.Description.String = strings.TrimSpace(p.Description.String)
	p.SmallPrint.String = strings.TrimSpace(p.SmallPrint.String)

	return validator.
		New("heading").
		Required().
		Validate(p.Heading)
}

// NewBaseProduct turns the input data into a BaseProduct.
func (p BaseProductInput) NewBaseProduct(creator string) BaseProduct {
	return BaseProduct{
		ID:               genProductID(),
		BaseProductInput: p,
		CreatedUTC:       chrono.TimeNow(),
		UpdatedUTC:       chrono.TimeNow(),
		CreatedBy:        creator,
	}
}

// BaseProduct defines a product without plans.
type BaseProduct struct {
	ID string `json:"id" db:"product_id"`
	BaseProductInput
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
	CreatedBy  string      `json:"createdByy" db:"created_by"`
}

// ProductInput defines the input data to create a product,
// with optional plans.
type ProductInput struct {
	BaseProductInput
	Plans []PlanInput `json:"plans"`
}

// Product turns input data into a Product.
func (p ProductInput) Product(creator string) Product {
	product := p.BaseProductInput.NewBaseProduct(creator)

	var plans = make([]BasePlan, 0)
	for _, v := range p.Plans {
		// Don't forget to add product id to plan.
		// Call NewPlan() won't add it since it assumes to
		// be provided by client.
		v.ProductID = product.ID
		plans = append(plans, v.NewPlan(creator))
	}

	return Product{
		BaseProduct: product,
		Plans:       plans,
	}
}

// Product is a product containing plans without discount
type Product struct {
	BaseProduct
	Plans []BasePlan `json:"plans"`
}

func (p Product) Update(input BaseProductInput) Product {
	p.Heading = input.Heading
	p.Description = input.Description
	p.SmallPrint = input.SmallPrint

	return p
}

// ProductSchema is used to hold db scan data for a list of product with plans retrieved as a JSON string.
type ProductSchema struct {
	BaseProduct
	Plans null.String `db:"plans"`
}

// ListedProduct converts data retrieve from db.
func (s ProductSchema) Product() (Product, error) {
	var plans = make([]BasePlan, 0)

	err := json.Unmarshal([]byte(s.Plans.String), &plans)
	if err != nil {
		return Product{}, err
	}

	return Product{
		BaseProduct: s.BaseProduct,
		Plans:       plans,
	}, nil
}
