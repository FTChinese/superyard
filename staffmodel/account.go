package staffmodel

type StaffAccount struct {
	ID          int    `json:"id"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
	Department  string `json:"department"`
	Groups      int    `json:"groups"`
	MyftID      string `json:"myftId"`
}
