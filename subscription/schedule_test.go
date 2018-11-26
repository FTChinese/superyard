package subscription

import (
	"testing"
)

func TestNewShedule(t *testing.T) {
	id, err := devEnv.NewSchedule(mockSchedule)

	if err != nil {
		t.Error(err)
	}

	t.Log(id)
}
