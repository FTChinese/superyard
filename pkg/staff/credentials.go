package staff

// Login specifies the the fields used for authentication
type Credentials struct {
	UserName string `db:"user_name"`
	Password string `db:"password"`
}
