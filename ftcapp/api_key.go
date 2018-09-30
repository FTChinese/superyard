package ftcapp

// APIKey contains data for an access token
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
