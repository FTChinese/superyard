package test

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/superyard/models/util"
	"log"
	"math/rand"
	"time"
)

var DBX *sqlx.DB
var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {

	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var dbConn util.Conn
	err = viper.UnmarshalKey("mysql.dev", &dbConn)
	if err != nil {
		log.Fatal(err)
	}

	DBX, err = util.NewDBX(dbConn)
	if err != nil {
		log.Fatal(err)
	}
}
