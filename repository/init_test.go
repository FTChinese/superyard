package repository

import (
	"database/sql"
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"github.com/Pallinder/go-randomdata"
	"github.com/guregu/null"
	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/backyard-api/models/util"
	"strings"
)

const (
	myFtcID    = "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae"
	myFtcEmail = "neefrankie@163.com"
	myUnionID  = "ogfvwjk6bFqv2yQpOrac0J3PqA0o"
)

var db *sql.DB

func init() {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	db, err = sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")
	if err != nil {
		panic(err)
	}
}

func genPassword() string {
	return fake.Password(8, 20, false, true, false)
}

func genUnionID() string {
	id, _ := gorest.RandomBase64(21)
	return id
}

func generateAvatarURL() string {
	return fmt.Sprintf("http://thirdwx.qlogo.cn/mmopen/vi_32/%s/132", fake.CharactersN(90))
}

func genOrderID() string {
	id, _ := gorest.RandomHex(8)

	return "FT" + strings.ToUpper(id)
}

func clientApp() util.ClientApp {
	return util.ClientApp{
		ClientType: enum.Platform(randomdata.Number(1, 4)),
		Version:    null.StringFrom("1.1.1"),
		UserIP:     null.StringFrom(fake.IPv4()),
		UserAgent:  null.StringFrom(fake.UserAgent()),
	}
}

func generateToken() string {
	token, _ := gorest.RandomBase64(82)
	return token
}

func generateWxID() string {
	id, _ := gorest.RandomBase64(21)
	return id
}
