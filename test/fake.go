package test

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/guregu/null"
	"github.com/icrowley/fake"
	"gitlab.com/ftchinese/backyard-api/android"
	"math/rand"
	"time"
)

// GenVersion creates a semantic version string.
func GenVersion() string {
	return fmt.Sprintf("%d.%d.%d", randomdata.Number(10), randomdata.Number(1, 10), randomdata.Number(1, 10))
}

func SemanticVersion() string {
	return "v" + GenVersion()
}

func Int64() int64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Int63n(10000)
}

func FakeURL() string {
	return fmt.Sprintf("https://www.%s/%s", fake.DomainName(), fake.Word())
}

func AndroidMock() android.Release {
	return android.Release{
		VersionName: SemanticVersion(),
		VersionCode: Int64(),
		Body:        null.StringFrom(fake.Paragraphs()),
		BinaryURL:   FakeURL(),
	}
}
