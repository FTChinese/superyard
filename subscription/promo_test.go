package subscription

import (
	"testing"
)

func TestRetrievePromo(t *testing.T) {
	d, err := devEnv.RetrievePromo(2)

	if err != nil {
		t.Error(err)
	}

	t.Log(d)
}

func TestListPromo(t *testing.T) {
	sch, err := devEnv.ListPromo(1, 10)

	if err != nil {
		t.Error(err)
	}

	t.Log(sch)
}
