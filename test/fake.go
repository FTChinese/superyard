package test

import (
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"github.com/Pallinder/go-randomdata"
	"github.com/guregu/null"
	"github.com/icrowley/fake"
	"gitlab.com/ftchinese/backyard-api/models/android"
	"gitlab.com/ftchinese/backyard-api/models/util"
)

func RandomClientApp() util.ClientApp {
	return util.ClientApp{
		ClientType: enum.Platform(randomdata.Number(1, 4)),
		Version:    null.StringFrom(GenVersion()),
		UserIP:     null.StringFrom(randomdata.IpV4Address()),
		UserAgent:  null.StringFrom(randomdata.UserAgentString()),
	}
}

// GenVersion creates a semantic version string.
func GenVersion() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		Rand.Intn(10),
		Rand.Intn(10),
		Rand.Intn(10),
	)
}

func SemanticVersion() string {
	return "v" + GenVersion()
}

func GetCusID() string {
	id, _ := gorest.RandomBase64(9)
	return "cus_" + id
}

func GenWxID() string {
	id, _ := gorest.RandomBase64(21)
	return id
}

func GenAvatar() string {
	var gender = []string{"men", "women"}

	n := randomdata.Number(1, 35)
	g := gender[randomdata.Number(0, 2)]

	return fmt.Sprintf("https://randomuser.me/api/portraits/thumb/%s/%d.jpg", g, n)
}

func FakeURL() string {
	return fmt.Sprintf("https://www.%s/%s", fake.DomainName(), fake.Word())
}

func AndroidMock() android.Release {
	return android.Release{
		VersionName: SemanticVersion(),
		VersionCode: Rand.Int63n(1000),
		Body:        null.StringFrom(fake.Paragraphs()),
		ApkURL:      FakeURL(),
	}
}
