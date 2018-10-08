package test

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	e := os.Environ()

	t.Log(e)
}
