package util

import (
	"encoding/json"
	"io"
)

// Parse parses input data to struct
func Parse(data io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(data)
	defer data.Close()

	return dec.Decode(v)
}
