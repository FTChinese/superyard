package fetch

import "testing"

func TestURLBuilder_String(t *testing.T) {
	b := NewURLBuilder("https://example.org").
		AddPath("demo").
		AddQuery("live", "true").
		SetRawQuery("page=1&per_page=20")

	t.Logf("%s", b.String())
}
