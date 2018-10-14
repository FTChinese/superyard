package staff

import "testing"

func TestCreateResetToken(t *testing.T) {
	token, err := devEnv.createResetToken(mockAccount.Email)

	if err != nil {
		t.Error(err)
	}

	t.Log(token)
}

func TestVerifyResetToken(t *testing.T) {
	token, err := devEnv.createResetToken(mockAccount.Email)

	a, err := devEnv.VerifyResetToken(token)

	if err != nil {
		t.Error(err)
	}

	t.Log(a)
}

func TestDeleteResetToken(t *testing.T) {
	token, err := devEnv.createResetToken(mockAccount.Email)

	t.Log(token)

	err = devEnv.deleteResetToken(token)

	if err != nil {
		t.Error(err)
	}
}
