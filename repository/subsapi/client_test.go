package subsapi

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestClient_RefreshPaywall(t *testing.T) {
	faker.MustConfigViper()

	c := NewClient(true)

	resp, err := c.RefreshPaywall()
	if err != nil {
		t.Error(err)
	}

	if assert.Equal(t, resp.StatusCode, 200) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		t.Logf("%s", body)
	}
}
