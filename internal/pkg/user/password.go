package user

import (
	"fmt"
	"time"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/pkg/conv"
)

type PwResetSession struct {
	Token      conv.HexBin `gorm:"primaryKey;column:token"`
	Email      string      `gorm:"column:email"`
	IsUsed     bool        `gorm:"column:is_used"`
	ExpiresIn  int64       `gorm:"column:expires_in"`
	CreatedUTC chrono.Time `gorm:"column:created_utc"`
}

func (PwResetSession) TableName() string {
	return "backyard.password_reset"
}

// NewPwResetSession creates a new PwResetSession instance
// based on request body which contains a required `email`
// field, and an optionally `sourceUrl` field.
func NewPwResetSession(email string) (PwResetSession, error) {
	token, err := conv.RandomHexBin(32)
	if err != nil {
		return PwResetSession{}, err
	}

	return PwResetSession{
		Token:      token,
		Email:      email,
		IsUsed:     false,
		ExpiresIn:  10800,
		CreatedUTC: chrono.TimeUTCNow(),
	}, nil
}

func MustNewPwResetSession(email string) PwResetSession {
	s, err := NewPwResetSession(email)

	if err != nil {
		panic(err)
	}

	return s
}

func (s PwResetSession) Disable() PwResetSession {
	s.IsUsed = true
	return s
}

func (s PwResetSession) BuildURL(baseURL string) string {
	if baseURL == "" {
		baseURL = "https://superyard.ftchinese.com/auth/forgot-password"
	}

	return fmt.Sprintf("%s/%s", baseURL, s.Token)
}

// IsExpired tests whether an existing PwResetSession is expired.
func (s PwResetSession) IsExpired() bool {
	return s.CreatedUTC.Add(time.Second * time.Duration(s.ExpiresIn)).Before(time.Now())
}
