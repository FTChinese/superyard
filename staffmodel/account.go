package staffmodel

// Account contains data returned after user authenticated successfully
type Account struct {
	ID          int    `json:"id"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
	Department  string `json:"department"`
	Groups      int    `json:"groups"`
	MyftID      string `json:"myftId"`
}
