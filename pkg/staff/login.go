package staff

// Login specifies the the fields used for authentication
type Login struct {
	UserName string `json:"userName" db:"user_name"`
	Password string `json:"password" db:"password"`
}
