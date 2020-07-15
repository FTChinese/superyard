package staff

import (
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
)

func GenStaffID() string {
	return "stf_" + rand.String(12)
}

type SignUp struct {
	Account
	Password string `json:"password" db:"password"`
}

func NewSignUp(input InputData) SignUp {
	input.ID = null.StringFrom(GenStaffID())

	return input.SignUp
}
