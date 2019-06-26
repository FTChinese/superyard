package test

import (
	"github.com/icrowley/fake"
	"testing"
)

func TestRandomURL(t *testing.T) {
	t.Log(fake.DomainName())
	t.Log(fake.DomainZone())
	t.Log(fake.TopLevelDomain())

	t.Log(fake.UserName())
	t.Log(fake.Characters())
	t.Log(fake.Word())
	t.Log(fake.Words())
}

func TestFakeUrl(t *testing.T) {
	t.Log(FakeURL())
}

func TestAndroidMock(t *testing.T) {
	t.Logf("Mock release: %+v", AndroidMock())
}
