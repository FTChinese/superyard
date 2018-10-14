package admin

import (
	"database/sql"

	"gitlab.com/ftchinese/backyard-api/staff"
)

func newDevEnv() Env {
	db, err := sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")

	if err != nil {
		panic(err)
	}

	return Env{DB: db}
}

var devEnv = newDevEnv()

var mockStaff = staff.Account{
	UserName:     "foo.bar",
	Email:        "foo.bar@ftchinese.com",
	DisplayName:  "Foo Bar",
	Department:   "tech",
	GroupMembers: 3,
}

var mockMyft = staff.MyftAccount{
	ID:    "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae",
	Email: "weiguo.ni@ftchinese.com",
}
