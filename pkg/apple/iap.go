package apple

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/pkg/paywall"
)

// Subscription contains a user's subscription data.
// It it built from app store's verification response.
// The original transaction id is used to uniquely identify a user.
type Subscription struct {
	Environment           Environment `json:"environment" db:"environment"`
	OriginalTransactionID string      `json:"originalTransactionId" db:"original_transaction_id"`
	LastTransactionID     string      `json:"lastTransactionId" db:"last_transaction_id"`
	ProductID             string      `json:"productId" db:"product_id"`
	PurchaseDateUTC       chrono.Time `json:"purchaseDateUtc" db:"purchase_date_utc"`
	ExpiresDateUTC        chrono.Time `json:"expiresDateUtc" db:"expires_date_utc"`
	paywall.Edition
	AutoRenewal bool        `json:"autoRenewal" db:"auto_renewal"`
	CreatedUTC  chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC  chrono.Time `json:"updatedUtc" db:"updated_utc"`
}
