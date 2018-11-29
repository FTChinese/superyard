package subscription

import (
	"testing"
)

func TestRetrievePromo(t *testing.T) {
	p, err := devEnv.RetrievePromo(2)

	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v\n", p)
}

func TestListPromo(t *testing.T) {
	promos, err := devEnv.ListPromo(1, 10)

	if err != nil {
		t.Error(err)
	}

	t.Log(promos)
}

func TestDisablePromo(t *testing.T) {
	err := devEnv.DisablePromo(1)

	if err != nil {
		t.Error(err)
	}
}
