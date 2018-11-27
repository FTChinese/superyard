package subscription

import "testing"

func TestSaveBanner(t *testing.T) {
	err := devEnv.SaveBanner(1, mockBanner)

	if err != nil {
		t.Error(err)
	}
}
