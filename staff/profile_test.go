package staff

import "testing"

func TestProfile(t *testing.T) {
	p, err := devEnv.Profile(mockAccount.UserName)

	if err != nil {
		t.Error(err)
	}

	t.Log(p)
}

func TestUpdateName(t *testing.T) {
	err := devEnv.UpdateName(mockAccount.UserName, mockAccount.DisplayName)

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateEmail(t *testing.T) {
	err := devEnv.UpdateEmail(mockAccount.UserName, mockAccount.Email)

	if err != nil {
		t.Error(err)
	}
}

func TestUpdatePassword(t *testing.T) {
	pass := Password{
		Old: "12345678",
		New: "12345678",
	}
	err := devEnv.UpdatePassword(mockAccount.UserName, pass)

	if err != nil {
		t.Error(err)
	}
}
