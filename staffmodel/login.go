package staffmodel

// StaffLogin specifies the the fields used for authentication
type StaffLogin struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	UserIP   string `json:"userIp"`
}
