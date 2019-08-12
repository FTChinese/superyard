package test

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/backyard-api/models/util"
)

var DB *sql.DB
var DBX *sqlx.DB

func init() {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var dbConn util.Conn
	err = viper.UnmarshalKey("mysql.dev", &dbConn)
	if err != nil {
		panic(err)
	}

	DB, err = util.NewDB(dbConn)
	if err != nil {
		panic(err)
	}

	DBX, err = util.NewDBX(dbConn)
	if err != nil {
		panic(err)
	}
}
