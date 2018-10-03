package ftcuser

// Membership contains a user's membership information
type Membership struct {
	Tier         string `json:"tier"`
	BillingCycle string `json:"billingCycle"`
	Start        string `json:"startAt"`
	Expire       string `json:"expireAt"`
}

// LoginHistory shows a user's login footprint
type LoginHistory struct {
	AuthType      string `json:"authType"`
	ClientType    string `json:"clientType"`
	ClientVersion string `json:"clientVersion"`
	UserIP        string `json:"userIp"`
	CreatedAt     string `json:"createdAt"`
}

// User is a registerd ftc user
type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Gender       string     `json:"gender"`
	FamilyName   string     `json:"familyName"`
	GivenName    string     `json:"givenName"`
	MobileNumber string     `json:"mobileNumber"`
	Birthdate    string     `json:"birthdate"`
	Address      string     `json:"address"`
	CreatedAt    string     `json:"createdAt"`
	Membership   Membership `json:"membership"`
}
