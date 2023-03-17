package user

import (
	"fmt"
	"time"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/pkg/conv"
)

// TODO: how to save Token as varbinary wity Gorm?
type PwResetSession struct {
	Token      conv.HexStr
	Email      string
	IsUsed     bool
	ExpiresIn  int64
	CreatedUTC chrono.Time
}

func (PwResetSession) TableName() string {
	return "backyard.password_reset"
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
		Token:      conv.HexStr(token),
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

type PasswordVerifier struct {
	StaffID     string
	OldPassword string
}

const StmtUpdatePassword = `
UPDATE cmstmp01.managers
	SET password = MD5(?)
WHERE username = ?
LIMIT 1`

const StmtDisableResetToken = `
UPDATE backyard.password_reset
SET is_used = 1
WHERE token = UNHEX(?)
LIMIT 1`
