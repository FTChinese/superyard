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

func TestEnablePromo(t *testing.T) {
	err := devEnv.EnablePromo(1, true)

	if err != nil {
		t.Error(err)
	}
}

func TestDisablePromo(t *testing.T) {
	err := devEnv.EnablePromo(1, false)

	if err != nil {
		t.Error(err)
	}
}
