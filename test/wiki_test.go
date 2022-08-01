//go:build !production

package test

import (
	"testing"
)

func TestNewArticle(t *testing.T) {
	t.Logf("Wiki article: %+v", NewArticle())
}
