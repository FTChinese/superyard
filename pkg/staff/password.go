package staff

import (
	"fmt"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"time"
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

func MustNewPwResetSession(email string) PwResetSession {
	s, err := NewPwResetSession(email)

	if err != nil {
		panic(err)
	}

	return s
}

func (s PwResetSession) BuildURL() string {
	return fmt.Sprintf("%s/%s", s.SourceURL, s.Token)
}

// IsExpired tests whether an existing PwResetSession is expired.
func (s PwResetSession) IsExpired() bool {
	return s.CreatedUTC.Add(time.Second * time.Duration(s.ExpiresIn)).Before(time.Now())
}

type PasswordVerifier struct {
	StaffID     string
	OldPassword string
}
