package ghapi

import (
	"github.com/FTChinese/superyard/faker"
	"testing"
)

func TestClient_GetAndroidLatestRelease(t *testing.T) {
	c := MustNewClient()

	r, respErr := c.GetAndroidLatestRelease()

	if respErr != nil {
		t.Error(respErr)
		return
	}

	t.Logf("%s", faker.MustMarshalIndent(r))
}

func TestClient_GetAndroidRelease(t *testing.T) {
	c := MustNewClient()

	r, respErr := c.GetAndroidRelease("v3.5.0")
	if respErr != nil {
		t.Error(respErr)
		return
	}

	t.Logf("%s", faker.MustMarshalIndent(r))
}

func TestClient_GetAndroidGradleFile(t *testing.T) {
	c := MustNewClient()

	content, respErr := c.GetAndroidGradleFile("v3.5.0")
	if respErr != nil {
		t.Error(respErr)
		return
	}

	t.Logf("%s", faker.MustMarshalIndent(content))
}
