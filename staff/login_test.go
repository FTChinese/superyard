package staff

import "testing"

func TestAuth(t *testing.T) {
	a, err := devEnv.Auth(mockLogin)

	if err != nil {
		t.Error(err)
	}

	t.Log(a)
}

func TestUpdateLoginHistory(t *testing.T) {
	err := devEnv.updateLoginHistory(mockLogin)

	if err != nil {
		t.Error(err)
	}
}
