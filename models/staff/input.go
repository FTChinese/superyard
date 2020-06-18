package staff

// Input is used as the JSON unmarshal target of a staff's data.
// Each request have different combination of those fields:
// Login: UserName + Password
// Request password reset: Email + SourceURL
// Reset password: Token + Password
// Set email: Email
// Update display name: DisplayName
// Update password: Password + OldPassword
type Input struct {
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	SourceURL   string `json:"sourceUrl"` // The URL to compose an email.
	Token       string `json:"token"`
	DisplayName string `json:"displayName"`
	OldPassword string `json:"oldPassword"`
}
