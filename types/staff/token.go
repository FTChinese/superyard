package staff

import "github.com/FTChinese/go-rest"

// TokenHolder holds a unique token for an email address.
// The token is readonly once generated.
type TokenHolder struct {
	t     string
	email string
}

// NewTokenHolder creates a new instance of TokenHolder with a 32 bytes random number.
func NewTokenHolder(email string) (TokenHolder, error) {
	token, err := gorest.RandomHex(32)

	if err != nil {
		return TokenHolder{}, err
	}

	return TokenHolder{t: token, email: email}, nil
}

// GetToken returns the generated token string.
func (h TokenHolder) GetToken() string {
	return h.t
}

// GetEmail returns the email associated with the token.
func (h TokenHolder) GetEmail() string {
	return h.email
}
