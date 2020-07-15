package config

import (
	"github.com/FTChinese/go-rest/connect"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Debug   bool
	Version string
	BuiltAt string
	Year    int
}

func GetConn(key string) (connect.Connect, error) {
	var conn connect.Connect
	err := viper.UnmarshalKey(key, &conn)
	if err != nil {
		return connect.Connect{}, err
	}

	return conn, nil
}

func (c Config) MustGetDBConn(key string) connect.Connect {
	var conn connect.Connect
	var err error

	if c.Debug {
		conn, err = GetConn("mysql.dev")
	} else {
		conn, err = GetConn(key)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using mysql server %s. Debugging: %t", conn.Host, c.Debug)

	return conn
}

func MustGetEmailConn() connect.Connect {

	conn, err := GetConn("email.ftc")
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func MustGetHanqiConn() connect.Connect {
	conn, err := GetConn("email.hanqi")
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
