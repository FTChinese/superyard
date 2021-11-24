package conv

import "strconv"

func ParseInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 0)
}
