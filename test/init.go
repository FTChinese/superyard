package test

import (
	"github.com/FTChinese/go-rest/connect"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/superyard/pkg/db"
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

	var dbConn connect.Connect
	err = viper.UnmarshalKey("mysql.dev", &dbConn)
	if err != nil {
		log.Fatal(err)
	}

	DBX, err = db.NewDB(dbConn)
	if err != nil {
		log.Fatal(err)
	}
}
