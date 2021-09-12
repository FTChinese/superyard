// +build !production

package faker

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

func MustMarshalIndent(v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "\t")

	if err != nil {
		panic(err)
	}

	return b
}

func MustReadBody(body io.Reader) []byte {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}

	return b
}
