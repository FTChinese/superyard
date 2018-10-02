package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/tidwall/gjson"
)

// Parse parses input data to struct
func Parse(data io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(data)
	defer data.Close()

	return dec.Decode(v)
}

// GetJSONString get a string field from http request body
func GetJSONString(data io.ReadCloser, path string) (string, error) {
	b, err := ioutil.ReadAll(data)
	defer data.Close()

	if err != nil {
		return "", ErrBadRequest
	}

	result := gjson.GetBytes(b, path)

	if !result.Exists() {
		ue := UnprocessableError{
			Field: path,
			Code:  CodeMissingField,
		}

		return "", ue
	}

	value := strings.TrimSpace(result.String())

	if value == "" {
		ue := UnprocessableError{
			Field: path,
			Code:  CodeMissingField,
		}

		return "", ue
	}

	return value, nil
}
