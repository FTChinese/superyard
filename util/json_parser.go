package util

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/tidwall/gjson"
)

// Flags telling caller how to respond to client errors
var (
	ErrNoField = errors.New("json: no field found for key")
)

// Parse parses input data to struct
func Parse(data io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(data)
	defer data.Close()

	return dec.Decode(v)
}

// GetJSONString get a string field from http request body
// Return empty string even if the passed in data does not contain the required key.
func GetJSONString(data io.ReadCloser, path string) (string, error) {
	b, err := ioutil.ReadAll(data)
	defer data.Close()

	if err != nil {
		return "", ErrBadRequest
	}

	result := gjson.GetBytes(b, path)

	if !result.Exists() {
		return "", nil
	}

	value := strings.TrimSpace(result.String())

	return value, nil
}
