package staff

import (
	"github.com/guregu/null"
	"github.com/icrowley/fake"
	"testing"
)

func mockAccount() Account {
	a, err := NewAccount()
	if err != nil {
		panic(err)
	}

	a.Email = fake.EmailAddress()
	a.UserName = fake.UserName()
	a.DisplayName = null.StringFrom(fake.FullName())
	a.Department = null.StringFrom("tech")
	a.GroupMembers = 3

	return a
}

func TestAccount_SignUpParcel(t *testing.T) {

	p, err := mockAccount().SignUpParcel()

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", p)
}

func TestAccount_PasswordResetParcel(t *testing.T) {
	a := mockAccount()
	th, _ := a.TokenHolder()

	p, err := a.PasswordResetParcel(th.GetToken())

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", p)
}
