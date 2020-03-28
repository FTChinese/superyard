package employee

import "github.com/FTChinese/go-rest/rand"

func GenStaffID() string {
	return "stf_" + rand.String(12)
}
