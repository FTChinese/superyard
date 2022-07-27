//go:build !production

package faker

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/rand"
	"github.com/brianvoe/gofakeit/v5"
	"time"
)

func SeedGoFake() {
	gofakeit.Seed(time.Now().UnixNano())
}

func GenStripeSubID() string {
	id, _ := rand.Base64(9)
	return "sub_" + id
}

func GenStripePlanID() string {
	return "plan_" + rand.String(14)
}

func RandNumericString() string {
	return rand.StringWithCharset(9, "0123456789")
}

func GenAppleSubID() string {
	return "1000000" + RandNumericString()
}

func GenWxID() string {
	id, _ := gorest.RandomBase64(21)
	return id
}

func GenToken32Bytes() string {
	token, _ := gorest.RandomHex(32)
	return token
}

func GenLicenceID() string {
	return "lic_" + rand.String(12)
}

func SimplePassword() string {
	return gofakeit.Password(true, false, true, false, false, 8)
}
