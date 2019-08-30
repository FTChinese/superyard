package builder

import (
	"errors"
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

func WhereStaffAccount(p SearchParam) (*Where, error) {
	if p.Name != "" {
		return NewWhere().Append("user_name", p.Name).Limit(1), nil
	}

	if p.Email != "" {
		return NewWhere().Append("email", p.Email).Limit(1), nil
	}

	return nil, errors.New("empty where clause")
}
