package ftcapi

import (
	"testing"

	"github.com/icrowley/fake"
)

func TestNewAPIKey(t *testing.T) {
	err := devEnv.NewAPIKey(APIKey{
		Description: fake.Sentence(),
		CreatedBy:   mockApp.OwnedBy,
	})

	if err != nil {
		t.Error(err)
	}
}

func TestPersonalAPIKeys(t *testing.T) {
	keys, err := devEnv.PersonalAPIKeys(mockApp.OwnedBy)

	if err != nil {
		t.Error(err)
	}

	t.Log(keys)
}

func TestRemovePersonalAccess(t *testing.T) {
	err := devEnv.RemovePersonalAccess(3, mockApp.OwnedBy)

	if err != nil {
		t.Error(err)
	}
}
