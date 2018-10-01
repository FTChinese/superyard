package staff

import (
	"github.com/parnurzeal/gorequest"
	"gitlab.com/ftchinese/backyard-api/util"
)

const resetLetterURL = "http://localhost:8001/backyard/password-reset"

// PasswordReset is used as marshal target when user tries to reset password via email
type PasswordReset struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// LetterAddress finds a user's name, email and display name by verifying its email.
func (env Env) LetterAddress(email string) (LetterAddress, error) {
	query := `
	SELECT username AS userName,
		email,
		display_name AS displayName
	FROM backyard.staff
	WHERE email = ?
	LIMIT 1`

	var a LetterAddress
	err := env.DB.QueryRow(query, email).Scan(
		&a.UserName,
		&a.Email,
		&a.DisplayName,
	)

	if err != nil {
		staffLogger.WithField("location", "Verify password reset email").Error(err)

		return a, err
	}

	return a, nil
}

// CreateResetToken send a password reset token to a user's email
func (env Env) CreateResetToken(a LetterAddress) error {
	token, err := util.RandomHex(32)

	staffLogger.Infof("Password reset token: %s\n", token)

	if err != nil {
		staffLogger.
			WithField("location", "Generate password reset token").
			Error(err)

		return err
	}

	query := `
	INSERT INTO backyard.password_reset
    SET token = UNHEX(?),
		email = ?`

	_, err = env.DB.Exec(query, token, a.Email)

	if err != nil {
		staffLogger.WithField("location", "Save password reset token").Error(err)
		return err
	}

	request := gorequest.New()

	_, _, errs := request.Post(resetLetterURL).
		Send(map[string]string{
			"userName": a.UserName,
			"token":    token,
			"address":  a.Email,
		}).
		End()

	if errs != nil {
		staffLogger.WithField("location", "Send password reset letter").Error(errs)

		return errs[0]
	}

	return nil
}

// VerifyResetToken finds the account associated with a password reset token
// If an account associated with a token if found, allow user to reset password.
func (env Env) VerifyResetToken(token string) (LetterAddress, error) {
	query := `
	SELECT s.username AS userName,
		s.email,
		s.display_name AS displayName,
	FROM backyard.password_reset AS p
		LEFT JOIN backyard.staff AS s
		ON p.email = s.email
    WHERE token = UNHEX(?)
      AND DATE_ADD(created_utc, INTERVAL expires_in SECOND) > UTC_TIMESTAMP()
	LIMIT 1`

	var a LetterAddress
	err := env.DB.QueryRow(query, token).Scan(
		&a.UserName,
		&a.Email,
		&a.DisplayName,
	)

	if err != nil {
		staffLogger.WithField("location", "Find email associated with a reset token").Error(err)

		return a, err
	}

	return a, nil
}

// ResetPassword allows user to reset password after clicked the password reset link in its email.
func (env Env) ResetPassword(r PasswordReset) error {
	// First check if the token is associated with an account
	addr, err := env.VerifyResetToken(r.Token)

	if err != nil {
		return err
	}

	// The account associated with a token is found. Chnage password.
	err = env.changePassword(addr.UserName, r.Password)

	if err != nil {
		return err
	}

	return nil
}

// DeleteResetToken deletes a password reset token after it is used.
func (env Env) DeleteResetToken(token string) error {
	query := `
	DELETE FROM backyard.password_reset
    WHERE token = UNHEX(?)
	LIMIT 1`

	_, err := env.DB.Exec(query, token)

	if err != nil {
		staffLogger.WithField("location", "Deleting a used password reset token").Error(err)

		return err
	}

	return nil
}
