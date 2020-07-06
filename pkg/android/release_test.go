package android

import (
	"encoding/json"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	data := []byte(`
	{
		"versionName": "v2.1.3",
		"versionCode": 17,
		"body": "This version fixed crashes when migrating to ViewModel",
		"apkUrl": "https://www.Avaveo.edu/temporibu"
	}`)

	var r Release

	err := json.Unmarshal(data, &r)

	if err != nil {
		t.Error(err)
	}

	t.Logf("Release: %+v", r)
}
