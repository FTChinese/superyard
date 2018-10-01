package staff

import (
	"github.com/parnurzeal/gorequest"
	"gitlab.com/ftchinese/backyard-api/util"
)

const resetLetterURL = "http://localhost:8001/backyard/password-reset"

// PasswordReset is used as marshal target when user tries to reset password via email
type PasswordReset struct {
	Token    string `json:"token"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateResetToken send a password reset token to a user's email
// Workflow:
// 1. When request comes, get email from request;
// 2. Use email to find a user's account. 404 if not found;
// 3. Call this function to create a reset token adn send a leeter.CreateResetToken
func (env Env) CreateResetToken(a Account) error {
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

// EmailForReset finds the email associated with a password reset token
func (env Env) EmailForReset(token string) (string, error) {
	query := `
	SELECT email AS email
    FROM backyard.password_reset
    WHERE token = UNHEX(?)
      AND DATE_ADD(created_utc, INTERVAL expires_in SECOND) > UTC_TIMESTAMP()
	LIMIT 1`

	var email string
	err := env.DB.QueryRow(query, token).Scan(&email)

	if err != nil {
		staffLogger.WithField("location", "Find email associated with a reset token").Error(err)

		return email, err
	}

	return email, nil
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
