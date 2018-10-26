package controller

import "testing"

func TestStatsDate(t *testing.T) {
	s, e, err := normalizeTimeRange("2018-10-01", "2018-10-26")

	if err != nil {
		t.Error(err)
	}

	t.Logf("Start %s, end %s\n", s, e)
}
