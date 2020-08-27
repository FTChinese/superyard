package config

import "github.com/FTChinese/go-rest/connect"

func MustGetEmailConn() connect.Connect {

	conn, err := GetConn("email.ftc")
	if err != nil {
		panic(err)
	}

	return conn
}

func MustGetHanqiConn() connect.Connect {
	conn, err := GetConn("email.hanqi")
	if err != nil {
		panic(err)
	}

	return conn
}
