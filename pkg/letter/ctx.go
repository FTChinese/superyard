package letter

type CtxSignUp struct {
	DisplayName string
	LoginName   string
	Password    string
	LoginURL    string
}

type CtxPasswordReset struct {
	DisplayName string
	URL         string
}

type CtxConfirmOrder struct {
	Name           string
	OrderCreatedAt string
	OrderID        string
	OrderAmount    string
	PayMethod      string
	OrderStartDate string
	OrderEndDate   string
	Tier           string
	ExpirationDate string
}
