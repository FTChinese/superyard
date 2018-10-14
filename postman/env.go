package postman

import (
	"fmt"

	"github.com/go-mail/mail"
	log "github.com/sirupsen/logrus"
)

const (
	baseURL       = "http://superyard.ftchinese.com"
	senderNoReply = "report@ftchinese.com"
)

var logger = log.WithField("package", "backyar-api.postman")

// Env wraps email server connection
type Env struct {
	Dialer *mail.Dialer
}

// Parcel is the item for a postmane to deliver.
type Parcel struct {
	Name     string // recipient name
	Address  string // recipient email address
	Password string // password for a new account
	Token    string // random string to identify the uniqueness of this parcel
}

// SendAccount sends user's account information to its email.
func (env Env) SendAccount(p Parcel) error {
	m := mail.NewMessage()

	m.SetHeader("From", senderNoReply)
	m.SetAddressHeader("To", "", "")
	m.SetHeader("Subject", "Welcome")
	m.SetBody("text/plain", fmt.Sprintf(`
Dear ${reqBody.userName},

Welcome to join FTC.

The following is your credentials to sign in to FTC Content Management System.

Login name: ${reqBody.userName}
Password: ${reqBody.password}

The password is an automatically generated random string. You're suggested to sign in the Content Management System and change it as soon as possible.

You can login via: http://superyard.ftchinese.com.

This email contains sensitive data. Do not leak it to anyone else.

Thanks,
FTC Dev Team`))

	if err := env.Dialer.DialAndSend(m); err != nil {
		logger.WithField("location", "Send account letter").Error(err)

		return err
	}

	return nil
}

// SendPasswordReset sedns a letter to reset password
func (env Env) SendPasswordReset(p Parcel) error {
	m := mail.NewMessage()

	m.SetHeader("From", senderNoReply)
	m.SetAddressHeader("To", "", "")
	m.SetHeader("Subject", "[FTC CMS]Reset Your Pasword")
	m.SetBody("text/plain", fmt.Sprintf(`
${reqBody.userName}

We heard that you lost your FTC CMS password. Sorry about that!

But don’t worry! You can use the following link to reset your password:

http://superyard.ftchinese.com/password-reset/${reqBody.token}

If you don’t use this link within 3 hours, it will expire. To get a new password reset link, visit http://superyard.ftchinese.com.

Thanks,
FTC Dev Team`))

	if err := env.Dialer.DialAndSend(m); err != nil {
		logger.WithField("location", "Send password reset letter").Error(err)

		return err
	}

	return nil
}
