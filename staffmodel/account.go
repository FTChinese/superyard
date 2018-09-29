package staffmodel

// Role(s) a staff can have
const (
	RoleRoot      = 1
	RoleDeveloper = 2
	RoleEditor    = 4
	RoleWheel     = 8
	RoleSales     = 16
	RoleMarketing = 32
	RoleMetting   = 64
)

// Account contains data returned after user authenticated successfully
type Account struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	UserName     string `json:"userName"`
	DisplayName  string `json:"displayName"`
	Department   string `json:"department"`
	GroupMembers int    `json:"groupMembers"`
	MyftID       string `json:"myftId"`
}
