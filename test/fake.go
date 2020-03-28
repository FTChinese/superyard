package test

import (
	"fmt"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/brianvoe/gofakeit/v4"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/util"
	"time"
)

func SimplePassword() string {
	gofakeit.Seed(time.Now().UnixNano())

	return gofakeit.Password(true, false, true, false, false, 8)
}

func GenWxID() string {
	id, _ := gorest.RandomBase64(21)
	return id
}

func GenDeviceToken() string {
	token, err := gorest.RandomHex(32)

	if err != nil {
		panic(err)
	}

	return token
}

func GenPwResetToken() string {
	t, err := gorest.RandomHex(32)
	if err != nil {
		panic(err)
	}

	return t
}

func GenVrfToken() string {
	t, err := gorest.RandomHex(32)
	if err != nil {
		panic(err)
	}

	return t
}

func GenSubID() string {
	id, _ := gorest.RandomBase64(9)
	return "sub_" + id
}

func GetCusID() string {
	id, _ := gorest.RandomBase64(9)
	return "cus_" + id
}

const charset = "0123456789"

func randNumericString() string {
	return rand.StringWithCharset(9, charset)
}

func GenAppleSubID() string {
	return "1000000" + randNumericString()
}

func RandomPaymentMethod() enum.PayMethod {
	return enum.PayMethod(Rand.Intn(5))
}

func RandomClientApp() util.ClientApp {
	return util.ClientApp{
		ClientType: enum.Platform(Rand.Intn(3) + 1),
		Version:    null.StringFrom(GenVersion()),
		UserIP:     null.StringFrom(gofakeit.IPv4Address()),
		UserAgent:  null.StringFrom(gofakeit.UserAgent()),
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
