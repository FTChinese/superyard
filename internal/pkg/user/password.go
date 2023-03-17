package user

import (
	"fmt"
	"time"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
)

// TODO: how to save Token as varbinary wity Gorm?
type PwResetSession struct {
	Token      string
	Email      string
	IsUsed     bool
	ExpiresIn  int64
	CreatedUTC chrono.Time
	SourceURL  string `gorm:"-"`
}

func (PwResetSession) Tablename() string {
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
		Token:      token,
		Email:      email,
		IsUsed:     false,
		ExpiresIn:  10800,
		CreatedUTC: chrono.TimeNow(),
		SourceURL:  "https://superyard.ftchinese.com/auth/forgot-password",
	}, nil
}

func MustNewPwResetSession(email string) PwResetSession {
	s, err := NewPwResetSession(email)

	if err != nil {
		panic(err)
	}

	return s
}

func (s PwResetSession) WithSourceURL(url string) PwResetSession {
	if url != "" {
		s.SourceURL = url
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

const StmtUpdatePassword = `
UPDATE cmstmp01.managers
	SET password = MD5(?)
WHERE username = ?
LIMIT 1`

const StmtInsertPwResetSession = `
INSERT INTO backyard.password_reset
SET token = UNHEX(:token),
	email = :email,
	created_utc = UTC_TIMESTAMP()`

const StmtPwResetSession = `
SELECT LOWER(HEX(token)) AS token,
	email,
	is_used,
	expires_in,
	created_utc
FROM backyard.password_reset
WHERE token = UNHEX(?)
LIMIT 1`

const StmtDisableResetToken = `
UPDATE backyard.password_reset
SET is_used = 1
WHERE token = UNHEX(?)
LIMIT 1`
