package ftcapi

// APIKey is an OAuth 2.0 access token used by an app or person to access ftc api
type APIKey struct {
	ID          int    `json:"id"`
	Token       string `json:"token"`
	Description string `json:"description"`
	MyftID      string `json:"myftId"`
	CreateAt    string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	LastUsed    string `json:"lastUsed"`
	CreatedBy   string `json:"createdBy"`
	OwnedByApp  string `json:"ownedByApp"`
}
