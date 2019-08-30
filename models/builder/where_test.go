package builder

import (
	"testing"
)

func TestWhereBuilder(t *testing.T) {
	w := NewWhere().Append("user_name", "michael").Append("email", "michael@example.org").Limit(1)

	t.Logf("Whese clause: %+v", w.Build())
	t.Logf("Where values: %+v", w.Values)
}
