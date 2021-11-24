package conv

import "strconv"

func DefaultTrue(str string) bool {
	t, err := strconv.ParseBool(str)
	if err != nil {
		return true
	}

	return t
}

func DefaultFalse(str string) bool {
	t, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}

	return t
}
