package builder

import "strings"

type InBuilder struct {
	bindVars []string
	args     []interface{}
}

func NewInBuilder(args ...interface{}) InBuilder {
	b := InBuilder{
		bindVars: []string{},
		args:     []interface{}{},
	}

	for _, v := range args {
		b.bindVars = append(b.bindVars, "?")
		b.args = append(b.args, v)
	}

	return b
}

// Append adds extra args to the argument list passed to
// a prepared statement.
// Use cases:
// Whe your statement looks like this:
//
// WHERE user_id IN (?)
//	ORDER BY created_utc DESC
//	LIMIT ? OFFSET ?
//
// Since variadic arguments must be the last one the param list,
// and the value for in is an variadic, you cannot add
// value for LIMIT and OFFSET. Append adds all the following
// args in a slice and can be passed to Exec() method.
func (b InBuilder) Append(args ...interface{}) InBuilder {
	for _, v := range args {
		b.args = append(b.args, v)
	}

	return b
}

func (b InBuilder) PlaceHolder() string {
	return strings.Join(b.bindVars, ",")
}

func (b InBuilder) Values() []interface{} {
	return b.args
}
