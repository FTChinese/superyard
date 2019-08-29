package test

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/backyard-api/models/util"
	"math/rand"
	"time"
)

var DBX *sqlx.DB
var Rand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

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

	DBX, err = util.NewDBX(dbConn)
	if err != nil {
		panic(err)
	}
}
