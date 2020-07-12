package staff

import (
	"fmt"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
)

type PwResetSession struct {
	Token      string      `db:"token"`
	Email      string      `db:"email"`
	IsUsed     bool        `db:"is_used"`
	ExpiresIn  int64       `db:"expires_in"`
	CreatedUTC chrono.Time `db:"created_utc"`
	SourceURL  string
}

// NewPwResetSession creates a new PwResetSession instance
// based on request body which contains a required `email`
// field, and an optionally `sourceUrl` field.
func NewPwResetSession(email string) (PwResetSession, error) {
	token, err := gorest.RandomHex(32)
	if err != nil {
		return PwResetSession{}, err
	}

	return PwResetSession{
		Token:      token,
		Email:      email,
		IsUsed:     false,
		ExpiresIn:  10800,
		CreatedUTC: chrono.TimeNow(),
		SourceURL:  "https://superyard.ftchinese.com/password-reset",
	}, nil
}

func (s PwResetSession) BuildURL() string {
	return fmt.Sprintf("%s/%s", s.SourceURL, s.Token)
}

type PasswordHolder struct {
	ID       string `db:"staff_id"`
	Password string `db:"password"`
}
