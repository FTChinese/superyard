package subscription

import (
	"encoding/json"
	"testing"
)

func TestStringifyPlans(t *testing.T) {
	p, err := json.Marshal(mockPricing)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s\n", p)
}

func TestStringifyBanner(t *testing.T) {
	b, err := json.Marshal(mockBanner)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s\n", b)
}

func TestNewPromo(t *testing.T) {
	id, err := createPromo()

	if err != nil {
		t.Error(err)
	}

	t.Log(id)
}

func TestRetrievePromo(t *testing.T) {
	id, err := createPromo()

	p, err := devEnv.RetrievePromo(id)

	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v\n", p)
}

func TestListPromo(t *testing.T) {
	promos, err := devEnv.ListPromo(1, 5)

	if err != nil {
		t.Error(err)
	}

	t.Log(promos)
}

func TestDisablePromo(t *testing.T) {
	id, err := createPromo()

	err = devEnv.DisablePromo(id)

	if err != nil {
		t.Error(err)
	}
}
