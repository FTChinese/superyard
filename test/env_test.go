package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func TestEnv(t *testing.T) {
	e := os.Environ()

	t.Log(e)
}

func TestViper(t *testing.T) {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.ftc")
	err := viper.ReadInConfig()
	if err != nil {
		t.Error(err)
	}
	viper.WatchConfig()

	dbConfig := viper.GetStringMapString("mysql")

	t.Log(dbConfig)

	host := dbConfig["host"]
	port := dbConfig["port"]
	user := dbConfig["user"]
	pass := dbConfig["pass"]

	cfg := &mysql.Config{
		User:                 user,
		Passwd:               pass,
		Net:                  "tcp",
		Addr:                 host + ":" + port,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		t.Error(err)
	}

	err = db.Ping()

	if err != nil {
		t.Error(err)
	}
}
