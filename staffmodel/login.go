package staffmodel

// Login specifies the the fields used for authentication
type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	UserIP   string `json:"userIp"`
}
