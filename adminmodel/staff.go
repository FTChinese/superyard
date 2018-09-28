package adminmodel

// Staff creates a new employee
type Staff struct {
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	DisplayName  string `json:"displayName"`
	Department   string `json:"department"`
	GroupMembers int    `json:"groupMembers"`
}
