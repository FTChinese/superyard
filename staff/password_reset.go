package staff

import (
	"fmt"
	"github.com/FTChinese/go-rest/postoffice"
	"strings"
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

// savePasswordToken send a password reset token to a user's email
func (env Env) savePasswordToken(h TokenHolder) error {

	query := `
	INSERT INTO backyard.password_reset
    SET token = UNHEX(?),
		email = ?`

	_, err := env.DB.Exec(query, h.GetToken(), h.GetEmail())

	if err != nil {
		logger.WithField("trace", "savePasswordToken").Error(err)
		return err
	}

	return nil
}

// CreatePwResetParcel checks if an email exists and send a password reset letter to it if exists.
func (env Env) CreatePwResetParcel(email string) (postoffice.Parcel, error) {
	// First try to find the user associated with this email
	// Error could be ErrNoRows
	a, err := env.findAccount(colEmail, email)
	if err != nil {
		return postoffice.Parcel{}, err
	}

	th, err := a.TokenHolder()
	if err != nil {
		logger.WithField("trace", "CreatePwResetParcel").Error(err)
		return postoffice.Parcel{}, err
	}

	err = env.savePasswordToken(th)

	// Internal server error
	if err != nil {
		return postoffice.Parcel{}, err
	}

	return a.PasswordResetParcel(th.GetToken())
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
