package test

import (
	"fmt"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/rand"
	"github.com/brianvoe/gofakeit/v5"
	"strings"
	"time"
)

// genVersion creates a semantic version string.
func genVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		rand.IntRange(1, 10),
		rand.IntRange(1, 10),
		rand.IntRange(1, 10))
}

func simplePassword() string {
	gofakeit.Seed(time.Now().UnixNano())

	return gofakeit.Password(true, false, true, false, false, 8)
}

func genWxID() string {
	id, _ := rand.Base64(21)
	return id
}

func genOrderID() string {
	return "FT" + strings.ToUpper(rand.String(16))
}

func getCustomerID() string {
	id, _ := rand.Base64(9)
	return "cus_" + id
}

func genStripeSubID() string {
	id, _ := rand.Base64(9)
	return "sub_" + id
}

func genStripePlanID() string {
	return "plan_" + rand.String(14)
}

func randNumericString() string {
	return rand.StringWithCharset(9, "0123456789")
}

func genAppleSubID() string {
	return "1000000" + randNumericString()
}

func genAvatar() string {
	var gender = []string{"men", "women"}

	n := rand.IntRange(1, 35)
	g := gender[rand.IntRange(0, 2)]

	return fmt.Sprintf("https://randomuser.me/api/portraits/thumb/%s/%d.jpg", g, n)
}

func mustGenToken() string {
	token, err := gorest.RandomHex(32)

	if err != nil {
		panic(err)
	}

	return token
}

func genBirthday() string {
	return fmt.Sprintf("%d-%d-%d", rand.IntRange(1900, 2020), rand.IntRange(1, 13), rand.IntRange(1, 31))
}

func genLicenceID() string {
	return "lic_" + rand.String(12)
}

func semanticVersion() string {
	return "v" + genVersion()
}
