package admin

import "testing"

func TestVIPRoster(t *testing.T) {
	vips, err := devEnv.VIPRoster()

	if err != nil {
		t.Error(err)
	}

	t.Log(vips)
}

func TestGrantVIP(t *testing.T) {
	err := devEnv.GrantVIP(mockMyft.ID)

	if err != nil {
		t.Error(err)
	}
}

func TestRevokeVIP(t *testing.T) {
	err := devEnv.RevokeVIP(mockMyft.ID)

	if err != nil {
		t.Error(err)
	}
}
