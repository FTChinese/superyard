package builder

import (
	"fmt"
	"strings"
)

type Where struct {
	cols   []string
	Values []interface{}
	limit  int
}

func NewWhere() *Where {
	return &Where{}
}

func (w *Where) Append(col, val string) *Where {
	w.cols = append(w.cols, fmt.Sprintf(`%s = ?`, col))
	w.Values = append(w.Values, val)

	return w
}

func (w *Where) Limit(n int) *Where {
	w.limit = n

	return w
}

// Build generates a SQL WHERE clause:
// WHERE user_name = ?
//	AND email = ?
//	AND ....
func (w *Where) Build() string {
	c := " WHERE " + strings.Join(w.cols, " AND ")

	if w.limit > 0 {
		c = c + fmt.Sprintf(" LIMIT %d", w.limit)
	}

	return c
}
