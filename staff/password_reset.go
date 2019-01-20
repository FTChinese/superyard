package staff

import (
	"fmt"
	"strings"

	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/postman"
)

// PasswordReset is used as marshal target when user tries to reset password via email
type PasswordReset struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Sanitize removes leading and trailing space of each field
func (r *PasswordReset) Sanitize() {
	r.Token = strings.TrimSpace(r.Token)
	r.Password = strings.TrimSpace(r.Password)
}

func (env Env) findAccount(col sqlCol, value string) (Account, error) {
	query := fmt.Sprintf(`
	%s
	WHERE %s = ?
		AND is_active = 1
	LIMIT 1`, stmtAccount, string(col))

	var a Account
	err := env.DB.QueryRow(query, value).Scan(
		&a.ID,
		&a.Email,
		&a.UserName,
		&a.DisplayName,
		&a.Department,
		&a.GroupMembers,
	)

	if err != nil {
		logger.WithField("location", "Find account by username or email").Error(err)

		return a, err
	}

	return a, nil
}

// createResetToken send a password reset token to a user's email
func (env Env) createResetToken(email string) (string, error) {
	token, err := gorest.RandomHex(32)

	if err != nil {
		logger.
			WithField("location", "Generate password reset token").
			Error(err)

		return "", err
	}

	query := `
	INSERT INTO backyard.password_reset
    SET token = UNHEX(?),
		email = ?`

	_, err = env.DB.Exec(query, token, email)

	if err != nil {
		logger.WithField("location", "Save password reset token").Error(err)
		return "", err
	}

	return token, nil
}

// RequestResetToken checks if an email exists and send a password reset letter to it if exists.
func (env Env) RequestResetToken(email string) (postman.Parcel, error) {
	// First try to find the user associated with this email
	// Error could be ErrNoRows
	a, err := env.findAccount(colEmail, email)
	if err != nil {
		return postman.Parcel{}, err
	}

	token, err := env.createResetToken(email)

	// Internal server error
	if err != nil {
		return postman.Parcel{}, err
	}

	return postman.Parcel{
		Name:    a.UserName,
		Address: a.Email,
		Token:   token,
	}, nil
}

// VerifyResetToken finds the account associated with a password reset token
// If an account associated with a token is found, allow user to reset password.
func (env Env) VerifyResetToken(token string) (Account, error) {
	query := `
	SELECT s.username AS userName,
		IFNULL(s.email, '') AS email,
		IFNULL(s.display_name, '') AS displayName
	FROM backyard.password_reset AS r
		LEFT JOIN backyard.staff AS s
		ON r.email = s.email
    WHERE r.token = UNHEX(?)
	  AND DATE_ADD(r.created_utc, INTERVAL r.expires_in SECOND) > UTC_TIMESTAMP()
	  AND s.is_active = 1
	LIMIT 1`

	var a Account
	err := env.DB.QueryRow(query, token).Scan(
		&a.UserName,
		&a.Email,
		&a.DisplayName,
	)

	if err != nil {
		logger.WithField("location", "Find email associated with a reset token").Error(err)

		return a, err
	}

	return a, nil
}

// ResetPassword allows user to reset password after clicked the password reset link in its email.
func (env Env) ResetPassword(r PasswordReset) error {
	// First try to find the account associated with the token
	a, err := env.VerifyResetToken(r.Token)

	// SqlNoRows
	if err != nil {
		return err
	}

	// The account associated with a token is found. Chnage password.
	err = env.changePassword(a.UserName, r.Password)

	if err != nil {
		return err
	}

	err = env.deleteResetToken(r.Token)

	if err != nil {
		return err
	}

	return nil
}

// DeleteResetToken deletes a password reset token after it was used.
func (env Env) deleteResetToken(token string) error {
	query := `
	DELETE FROM backyard.password_reset
    WHERE token = UNHEX(?)
	LIMIT 1`

	_, err := env.DB.Exec(query, token)

	if err != nil {
		logger.WithField("location", "Deleting a used password reset token").Error(err)

		return err
	}

	return nil
}
