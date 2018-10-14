package ftcapi

import "database/sql"

func newDevEnv() Env {
	db, err := sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")

	if err != nil {
		panic(err)
	}

	return Env{DB: db}
}

var devEnv = newDevEnv()

var mockApp = App{
	Name:         "Next User",
	Slug:         "next-user",
	ClientID:     "88736db88fc0a4d689e1",
	ClientSecret: "***REMOVED***",
	RepoURL:      "https://github.com/FTChinese/next-user",
	Description:  "FTC user login, signup and settings",
	HomeURL:      "http://next.ftchinese.com/user",
	OwnedBy:      "foo.bar",
}
