package util

import (
	"testing"
	"time"
)

func TestDatetime(t *testing.T) {
	parsedTime, err := time.Parse(iso9075, "2012-08-24 14:32:38")

	if err != nil {
		t.Error(err)
	}

	result := parsedTime.UTC().Format(time.RFC3339)

	t.Log(result)
}

func TestUnix(t *testing.T) {
	result := FormatUnix(time.Now().Unix())

	t.Log(result)
}
