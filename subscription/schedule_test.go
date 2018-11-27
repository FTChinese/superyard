package subscription

import (
	"testing"
)

func TestNewShedule(t *testing.T) {
	id, err := devEnv.NewSchedule(mockSchedule, "weiguo.ni")

	if err != nil {
		t.Error(err)
	}

	t.Log(id)
}
