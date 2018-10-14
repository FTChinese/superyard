package staff

import (
	"database/sql"
	"testing"
)

func newDevEnv() Env {
	db, err := sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")

	if err != nil {
		panic(err)
	}

	return Env{DB: db}
}

var devEnv = newDevEnv()

var mockAccount = Account{
	UserName:     "foo.bar",
	Email:        "foo.bar@ftchinese.com",
	DisplayName:  "foo.bar",
	Department:   "tech",
	GroupMembers: 3,
}

var mockLogin = Login{
	UserName: "foo.bar",
	Password: "12345678",
	UserIP:   "127.0.0.1",
}

var mockMyft = MyftAccount{
	ID:    "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae",
	Email: "weiguo.ni@ftchinese.com",
}

const (
	mockMyftPass = "12345678"
)

func TestIsPasswordMatched(t *testing.T) {
	ok, err := devEnv.isPasswordMatched(mockLogin.UserName, mockLogin.Password)

	if err != nil {
		t.Error(err)
	}

	t.Log(ok)
}
