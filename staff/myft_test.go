package staff

import "testing"

func TestAuthMyft(t *testing.T) {
	a, err := devEnv.authMyft(MyftCredential{
		Email:    mockMyft.Email,
		Password: mockMyftPass,
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(a)
}

func TestSaveMyft(t *testing.T) {
	err := devEnv.saveMyft(mockAccount.UserName, mockMyft)

	if err != nil {
		t.Error(err)
	}
}

func TestListMyft(t *testing.T) {
	myfts, err := devEnv.ListMyft(mockAccount.UserName)

	if err != nil {
		t.Error(err)
	}

	t.Log(myfts)
}

func TestDeleteMyft(t *testing.T) {
	err := devEnv.DeleteMyft(mockAccount.UserName, mockMyft.ID)

	if err != nil {
		t.Error(err)
	}
}
