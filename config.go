package main

import (
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/superyard/models/util"
	"os"
)

type Config struct {
	Debug   bool
	Version string
	BuiltAt string
	Year    int
}

func (c Config) MustGetDBConn(key string) util.Conn {
	var conn util.Conn
	var err error

	if c.Debug {
		err = viper.UnmarshalKey("mysql.dev", &conn)
	} else {
		err = viper.UnmarshalKey(key, &conn)
	}

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	logger.Infof("Using mysql server %s. Debugging: %t", conn.Host, c.Debug)

	return conn
}

func MustGetEmailConn() util.Conn {
	var emailConn util.Conn
	err := viper.UnmarshalKey("email.ftc", &emailConn)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	return emailConn
}
