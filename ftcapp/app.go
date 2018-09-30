package ftcapp

// App contains information about an app used on ftchinese.com
type App struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RepoURL      string `json:"repoUrl"`
	Description  string `json:"description"`
	HomeURL      string `json:"homeUrl"`
	IsActive     bool   `json:"isActive"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	OwnedBy      string `json:"ownedBy"`
}

// Ownership contains data used to transfer an app's ownership
type Ownership struct {
	ID       int    `json:"id"`
	NewOwner string `json:"newOwner"`
	OldOwner string `json:"oldOwner"`
}
