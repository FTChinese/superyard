package staffmodel

// Password marshals request data for updating password
type Password struct {
	UserName string `json:"userName"`
	Old      string `json:"oldPassword"`
	New      string `json:"newPassword"`
}
