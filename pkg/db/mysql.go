package db

import (
	"fmt"
	"time"

	"github.com/FTChinese/superyard/pkg/config"
	"github.com/jmoiron/sqlx"

	"github.com/go-sql-driver/mysql"
)

func buildDSN(c config.Connect) string {
	cfg := &mysql.Config{
		User:   c.User,
		Passwd: c.Pass,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		// Always use UTC time.
		// Pay attention to how string values are specified.
		// The string value provided to MySQL must be quoted in single quote for this driver to work,
		// which means the single quote itself must be included in the string value.
		// The resulting string value passed to MySQL should look like: `%27<you string value>%27`
		// See ASCII Encoding Reference https://www.w3schools.com/tags/ref_urlencode.asp
		Params: map[string]string{
			"time_zone": `'+00:00'`,
		},
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
	}

	return cfg.FormatDSN()
}

func NewMyDB(c config.Connect) (*sqlx.DB, error) {
	cfg := &mysql.Config{
		User:   c.User,
		Passwd: c.Pass,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		// Always use UTC time.
		// Pay attention to how string values are specified.
		// The string value provided to MySQL must be quoted in single quote for this driver to work,
		// which means the single quote itself must be included in the string value.
		// The resulting string value passed to MySQL should look like: `%27<you string value>%27`
		// See ASCII Encoding Reference https://www.w3schools.com/tags/ref_urlencode.asp
		Params: map[string]string{
			"time_zone": `'+00:00'`,
		},
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
	}

	db, err := sqlx.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// When connecting to production server it throws error:
	// packets.go:36: unexpected EOF
	//
	// See https://github.com/go-sql-driver/mysql/issues/674
	db.SetConnMaxLifetime(time.Second)
	return db, nil
}

func MustNewMySQL(c config.Connect) *sqlx.DB {
	db, err := NewMyDB(c)
	if err != nil {
		panic(err)
	}

	return db
}

type ReadWriteMyDBs struct {
	Read   *sqlx.DB
	Write  *sqlx.DB
	Delete *sqlx.DB
}

func MustNewMyDBs(prod bool) ReadWriteMyDBs {
	return ReadWriteMyDBs{
		Read:   MustNewMySQL(config.MustMySQLReadConn(prod)),
		Write:  MustNewMySQL(config.MustMySQLWriteConn(prod)),
		Delete: MustNewMySQL(config.MustMySQLDeleteConn(prod)),
	}
}
