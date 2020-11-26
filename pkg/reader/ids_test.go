package reader

import (
	"github.com/gorilla/schema"
	"net/url"
	"testing"
)

func TestDecodeIDs(t *testing.T) {
	var decoder = schema.NewDecoder()

	data := url.Values{}
	data.Add("ftc_id", "abc")
	data.Add("union_id", "def")

	var id IDs
	err := decoder.Decode(&id, data)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", id)
}
