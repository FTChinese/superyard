package config

import (
	"log"

	"github.com/spf13/viper"
)

// Connect represents a connection to a server or database.
type Connect struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

func GetConn(key string) (Connect, error) {
	var conn Connect
	err := viper.UnmarshalKey(key, &conn)
	if err != nil {
		return Connect{}, err
	}

	return conn, nil
}

func MustMySQLConn(key string, prod bool) Connect {
	var conn Connect
	var err error

	conn, err = GetConn(key)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using mysql server %s. Production: %t", conn.Host, prod)

	return conn
}

func MustMySQLReadConn(prod bool) Connect {
	return MustMySQLConn("mysql.read", prod)
}

func MustMySQLWriteConn(prod bool) Connect {
	return MustMySQLConn("mysql.write", prod)
}

func MustMySQLDeleteConn(prod bool) Connect {
	return MustMySQLConn("mysql.delete", prod)
}

func MustGetEmailConn() Connect {

	conn, err := GetConn("email.ftc")
	if err != nil {
		panic(err)
	}

	return conn
}

func MustGetHanqiConn() Connect {
	conn, err := GetConn("email.hanqi")
	if err != nil {
		panic(err)
	}

	return conn
}
