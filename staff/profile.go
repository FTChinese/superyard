package staff

// Profile contains the full data of a staff
type Profile struct {
	ID            int    `json:"id"`
	UserName      string `json:"userName"`
	Email         string `json:"email"`
	IsActive      bool   `json:"isActive"`
	DisplayName   string `json:"displayName"`
	Department    string `json:"department"`
	GroupMembers  int    `json:"groupMembers"`
	MyftID        string `json:"myftId"`
	MyftEmail     string `json:"myftEmail"`
	CreatedAt     string `json:"createdAt"`
	DeactiviateAt string `json:"deactivatedAt"`
	UpdatedAt     string `json:"updatedAt"`
	LastLoginAt   string `json:"lastLoginAt"`
	LastLoginIP   string `json:"lastLoginIp"`
}
